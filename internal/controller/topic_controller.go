/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/pubsub"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/clock"
	ctrl "sigs.k8s.io/controller-runtime"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
)

// TopicReconciler reconciles a Topic object
type TopicReconciler struct {
	crclient.Client
	Scheme    *runtime.Scheme
	NewClient newPubSubClientFunc
	Recorder  record.EventRecorder
	Clock     clock.PassiveClock
}

const topicFinalizerName = "topic.googlecloudpubsuboperator.quipper.github.io/finalizer"

const topicCreationWarningDeadline = 10 * time.Minute

//+kubebuilder:rbac:groups=googlecloudpubsuboperator.quipper.github.io,resources=topics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=googlecloudpubsuboperator.quipper.github.io,resources=topics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=googlecloudpubsuboperator.quipper.github.io,resources=topics/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch

func (r *TopicReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var topic googlecloudpubsuboperatorv1.Topic
	if err := r.Client.Get(ctx, req.NamespacedName, &topic); err != nil {
		logger.Error(err, "unable to get the resource")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, crclient.IgnoreNotFound(err)
	}
	logger.Info("Found Topic resource")

	// examine DeletionTimestamp to determine if object is under deletion
	if topic.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&topic, topicFinalizerName) {
			controllerutil.AddFinalizer(&topic, topicFinalizerName)
			if err := r.Update(ctx, &topic); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&topic, topicFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteTopic(ctx, topic.Spec.ProjectID, topic.Spec.TopicID); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&topic, topicFinalizerName)
			if err := r.Update(ctx, &topic); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	t, err := r.createTopic(ctx, topic.Spec.ProjectID, topic.Spec.TopicID)
	if err != nil {
		if isPubSubAlreadyExistsError(err) {
			// don't treat as error
			logger.Info("Topic already exists in Cloud Pub/Sub", "error", err)

			// update iam policy
			if err = r.updateIamPolicy(ctx, topic.Spec.ProjectID, topic.Spec.TopicID, topic.Spec.Bindings); err != nil {
				logger.Error(err, "Unable to update iam policy for topic")
				return ctrl.Result{}, err
			}
			logger.Info("Updated IAM policy for topic", "id", topic.Spec.TopicID)

			topicPatch := crclient.MergeFrom(topic.DeepCopy())
			topic.Status.Phase = googlecloudpubsuboperatorv1.TopicStatusPhaseActive
			topic.Status.Message = ""
			if err := r.Client.Status().Patch(ctx, &topic, topicPatch); err != nil {
				logger.Error(err, "unable to update status")
				return ctrl.Result{}, err
			}
			logger.Info("Topic status has been patched to Active")
			return ctrl.Result{}, nil
		}

		if r.Clock.Since(topic.CreationTimestamp.Time) > topicCreationWarningDeadline {
			r.Recorder.Event(&topic, corev1.EventTypeWarning, "TopicCreateErrorDeadlineExceeded",
				fmt.Sprintf("Failed to create Topic for %s in Pub/Sub: %s", topicCreationWarningDeadline, err))
		} else {
			r.Recorder.Event(&topic, corev1.EventTypeNormal, "TopicCreateError",
				fmt.Sprintf("Failed to create Topic in Pub/Sub: %s", err))
		}
		topicPatch := crclient.MergeFrom(topic.DeepCopy())
		topic.Status.Phase = googlecloudpubsuboperatorv1.TopicStatusPhaseError
		topic.Status.Message = fmt.Sprintf("Pub/Sub error: %s", err)
		if err := r.Client.Status().Patch(ctx, &topic, topicPatch); err != nil {
			logger.Error(err, "unable to update status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}
	logger.Info("Created topic into Cloud Pub/Sub", "id", t.ID())

	// update iam policy
	if err = r.updateIamPolicy(ctx, topic.Spec.ProjectID, topic.Spec.TopicID, topic.Spec.Bindings); err != nil {
		logger.Error(err, "Unable to update iam policy for topic")
		return ctrl.Result{}, err
	}
	logger.Info("Updated IAM policy for topic", "id", t.ID())

	topicPatch := crclient.MergeFrom(topic.DeepCopy())
	topic.Status.Phase = googlecloudpubsuboperatorv1.TopicStatusPhaseActive
	topic.Status.Message = ""
	if err := r.Client.Status().Patch(ctx, &topic, topicPatch); err != nil {
		logger.Error(err, "unable to update status")
		return ctrl.Result{}, err
	}
	logger.Info("Topic status has been patched to Active")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TopicReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&googlecloudpubsuboperatorv1.Topic{}).
		Complete(r)
}

func (r *TopicReconciler) createTopic(ctx context.Context, projectID, topicID string) (*pubsub.Topic, error) {
	c, err := r.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer c.Close()

	t, err := c.CreateTopic(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("CreateTopic: %w", err)
	}

	return t, nil
}

func (r *TopicReconciler) deleteTopic(ctx context.Context, projectID, topicID string) error {
	c, err := r.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer c.Close()

	if err := c.Topic(topicID).Delete(ctx); err != nil {
		if isPubSubNotFoundError(err) {
			// for idempotent
			return nil
		}
		return fmt.Errorf("unable to delete topic %s: %w", topicID, err)
	}
	return nil
}

// updateIamPolicy sets the topic iam policy to match the defined spec, if the bindings are non nil.
func (r *TopicReconciler) updateIamPolicy(ctx context.Context, projectID, topicID string, bindings []googlecloudpubsuboperatorv1.IamBinding) error {
	if bindings != nil {
		c, err := r.NewClient(ctx, projectID)
		if err != nil {
			return fmt.Errorf("pubsub.NewClient: %w", err)
		}
		defer c.Close()

		policy := &iam.Policy{}
		for _, bind := range bindings {
			for _, member := range bind.ServiceAccounts {
				policy.Add(member, iam.RoleName(bind.Role))
			}
		}

		return c.Topic(topicID).IAM().SetPolicy(ctx, policy)
	}

	return nil
}
