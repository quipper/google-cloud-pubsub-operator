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
	"sigs.k8s.io/controller-runtime/pkg/log"

	pubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
)

// GoogleCloudPubSubTopicReconciler reconciles a GoogleCloudPubSubTopic object
type GoogleCloudPubSubTopicReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pubsuboperator.quipper.github.io,resources=googlecloudpubsubtopics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pubsuboperator.quipper.github.io,resources=googlecloudpubsubtopics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pubsuboperator.quipper.github.io,resources=googlecloudpubsubtopics/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GoogleCloudPubSubTopic object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *GoogleCloudPubSubTopicReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var topic pubsuboperatorv1.GoogleCloudPubSubTopic
	if err := r.Client.Get(ctx, req.NamespacedName, &topic); err != nil {
		logger.Error(err, "unable to get the resource")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Found the topic", "topic", topic)

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
func (r *GoogleCloudPubSubTopicReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pubsuboperatorv1.GoogleCloudPubSubTopic{}).
		Complete(r)
}

func createTopic(ctx context.Context, projectID, topicID string) (*pubsub.Topic, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()

	t, err := client.CreateTopic(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("CreateTopic: %w", err)
	}

	return t, nil
}
