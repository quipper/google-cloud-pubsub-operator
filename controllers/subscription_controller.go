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

	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
)

// SubscriptionReconciler reconciles a Subscription object
type SubscriptionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=googlecloudpubsuboperator.quipper.github.io,resources=subscriptions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=googlecloudpubsuboperator.quipper.github.io,resources=subscriptions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=googlecloudpubsuboperator.quipper.github.io,resources=subscriptions/finalizers,verbs=update

func (r *SubscriptionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var subscription googlecloudpubsuboperatorv1.Subscription
	if err := r.Client.Get(ctx, req.NamespacedName, &subscription); err != nil {
		logger.Error(err, "unable to get the resource")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Found the subscription", "subscription", subscription)

	s, err := createSubscription(ctx, subscription)
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

// SetupWithManager sets up the controller with the Manager.
func (r *SubscriptionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&googlecloudpubsuboperatorv1.Subscription{}).
		Complete(r)
}

func createSubscription(ctx context.Context, subscription googlecloudpubsuboperatorv1.Subscription) (*pubsub.Subscription, error) {
	c, err := pubsub.NewClient(ctx, subscription.Spec.SubscriptionProjectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer c.Close()

	topic := c.TopicInProject(subscription.Spec.TopicID, subscription.Spec.TopicProjectID)
	s, err := c.CreateSubscription(ctx, subscription.Spec.SubscriptionID, pubsub.SubscriptionConfig{
		Topic:            topic,
		ExpirationPolicy: 24 * time.Hour,
	})
	if err != nil {
		return nil, fmt.Errorf("CreateSubscription: %w", err)
	}

	return s, nil
}
