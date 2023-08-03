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

	"cloud.google.com/go/pubsub"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
)

const subscriptionFinalizerName = "subscription.googlecloudpubsuboperator.quipper.github.io/finalizer"

// SubscriptionReconciler reconciles a Subscription object
type SubscriptionReconciler struct {
	crclient.Client
	Scheme    *runtime.Scheme
	NewClient newPubSubClientFunc
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
		return ctrl.Result{}, crclient.IgnoreNotFound(err)
	}

	logger.Info("Found the subscription", "subscription", subscription)

	// examine DeletionTimestamp to determine if object is under deletion
	if subscription.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&subscription, subscriptionFinalizerName) {
			controllerutil.AddFinalizer(&subscription, subscriptionFinalizerName)
			if err := r.Update(ctx, &subscription); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&subscription, subscriptionFinalizerName) {
			// our finalizer is present, so lets handle any external dependency
			if err := r.deleteSubscription(ctx, subscription); err != nil {
				// if fail to delete the external dependency here, return with error
				// so that it can be retried
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(&subscription, subscriptionFinalizerName)
			if err := r.Update(ctx, &subscription); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	s, err := r.createSubscription(ctx, subscription)
	if err != nil {
		if isPubSubAlreadyExistsError(err) {
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

func (r *SubscriptionReconciler) createSubscription(ctx context.Context, subscription googlecloudpubsuboperatorv1.Subscription) (*pubsub.Subscription, error) {
	c, err := r.NewClient(ctx,
		subscription.Spec.SubscriptionProjectID)
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

func (r *SubscriptionReconciler) deleteSubscription(ctx context.Context, subscription googlecloudpubsuboperatorv1.Subscription) error {
	c, err := r.NewClient(ctx, subscription.Spec.SubscriptionProjectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer c.Close()

	err = c.Subscription(subscription.Spec.SubscriptionID).Delete(ctx)
	if err != nil {
		if isPubSubNotFoundError(err) {
			// for idempotent
			return nil
		}
		return fmt.Errorf("unable to delete subscription: %w", err)
	}
	return nil
}
