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

package controllers

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
)

// TopicReconciler reconciles a Topic object
type TopicReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const topicFinalizerName = "topic.googlecloudpubsuboperator.quipper.github.io/finalizer"

//+kubebuilder:rbac:groups=googlecloudpubsuboperator.quipper.github.io,resources=topics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=googlecloudpubsuboperator.quipper.github.io,resources=topics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=googlecloudpubsuboperator.quipper.github.io,resources=topics/finalizers,verbs=update

func (r *TopicReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var topic googlecloudpubsuboperatorv1.Topic
	if err := r.Client.Get(ctx, req.NamespacedName, &topic); err != nil {
		logger.Error(err, "unable to get the resource")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	logger.Info("Found the topic", "topic", topic)

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
			if err := deleteTopic(ctx, topic.Spec.ProjectID, topic.Spec.TopicID); err != nil {
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

	t, err := createTopic(ctx, topic.Spec.ProjectID, topic.Spec.TopicID)
	if err != nil {
		if gs, ok := gRPCStatusFromError(err); ok && gs.Code() == codes.AlreadyExists {
			// don't treat as error
			logger.Info("PubSub topic already exists")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	logger.Info(fmt.Sprintf("Topic created: %v", t.ID()), "topic", topic)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TopicReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&googlecloudpubsuboperatorv1.Topic{}).
		Complete(r)
}

func createTopic(ctx context.Context, projectID, topicID string) (*pubsub.Topic, error) {
	c, err := pubsub.NewClient(ctx, projectID)
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

func deleteTopic(ctx context.Context, projectID, topicID string) error {
	c, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer c.Close()

	if err := c.Topic(topicID).Delete(ctx); err != nil {
		if gs, ok := gRPCStatusFromError(err); ok && gs.Code() == codes.NotFound {
			// for idempotent
			return nil
		}
		return fmt.Errorf("unable to delete topic %s: %w", topicID, err)
	}
	return nil
}
