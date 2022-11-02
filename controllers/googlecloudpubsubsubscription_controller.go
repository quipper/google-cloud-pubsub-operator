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
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	pubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
)

// GoogleCloudPubSubSubscriptionReconciler reconciles a GoogleCloudPubSubSubscription object
type GoogleCloudPubSubSubscriptionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pubsuboperator.quipper.github.io,resources=googlecloudpubsubsubscriptions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pubsuboperator.quipper.github.io,resources=googlecloudpubsubsubscriptions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pubsuboperator.quipper.github.io,resources=googlecloudpubsubsubscriptions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GoogleCloudPubSubSubscription object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *GoogleCloudPubSubSubscriptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var subscription pubsuboperatorv1.GoogleCloudPubSubSubscription
	if err := r.Client.Get(ctx, req.NamespacedName, &subscription); err != nil {
		logger.Error(err, "unable to get the resource")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Found the subscription", "subscription", subscription)

	s, err := createSubscription(ctx, subscription.Spec.ProjectID, subscription.Spec.TopicID)
	if err != nil {
		if gs, ok := gRPCStatusFromError(err); ok && gs.Code() == codes.AlreadyExists {
			// don't treat as error
			logger.Info("PubSub subscription already exists")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	logger.Info(fmt.Sprintf("Subscription created: %v", s.ID()), "subscription", subscription)

	return ctrl.Result{}, nil
}

func createSubscription(ctx context.Context, projectID, topicID string) (*pubsub.Subscription, error) {
	c, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer c.Close()

	s, err := c.CreateSubscription(ctx, topicID, pubsub.SubscriptionConfig{
		ExpirationPolicy: 24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateSubscription: %w", err)
	}

	return s, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GoogleCloudPubSubSubscriptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pubsuboperatorv1.GoogleCloudPubSubSubscription{}).
		Complete(r)
}
