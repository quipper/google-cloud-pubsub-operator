package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
	"github.com/quipper/google-cloud-pubsub-operator/internal/pubsubtest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	//+kubebuilder:scaffold:imports
)

var _ = Describe("Subscription controller", func() {
	Context("When creating a Subscription resource", func() {
		It("Should create a Pub/Sub Subscription", func(ctx context.Context) {
			const projectID = "subscription-project-1"
			psClient, err := pubsubtest.NewClient(ctx, projectID, psServer)
			Expect(err).ShouldNot(HaveOccurred())

			By("Creating a Topic")
			topicID := "my-topic"
			_, err = psClient.CreateTopic(ctx, topicID)
			Expect(err).ShouldNot(HaveOccurred())

			By("Creating a Subscription")
			subscription := &googlecloudpubsuboperatorv1.Subscription{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "googlecloudpubsuboperator.quipper.github.io/v1",
					Kind:       "Subscription",
				},
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "example-",
					Namespace:    "default",
				},
				Spec: googlecloudpubsuboperatorv1.SubscriptionSpec{
					SubscriptionProjectID: projectID,
					SubscriptionID:        "my-subscription",
					TopicProjectID:        projectID,
					TopicID:               topicID,
				},
			}
			Expect(k8sClient.Create(ctx, subscription)).Should(Succeed())

			By("Checking if the Subscription exists")
			Eventually(func(g Gomega) {
				subscriptionExists, err := psClient.Subscription("my-subscription").Exists(ctx)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(subscriptionExists).Should(BeTrue())
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())
		})

		It("Should update the status if error", func(ctx context.Context) {
			const projectID = "subscription-project-2"

			By("Creating a Subscription")
			subscription := &googlecloudpubsuboperatorv1.Subscription{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "googlecloudpubsuboperator.quipper.github.io/v1",
					Kind:       "Subscription",
				},
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "example-",
					Namespace:    "default",
				},
				Spec: googlecloudpubsuboperatorv1.SubscriptionSpec{
					SubscriptionProjectID: projectID,
					SubscriptionID:        "my-subscription",
					// CreateSubscription API should fail because the topic does not exist.
					// We don't need to explicitly inject an error.
					TopicProjectID: projectID,
					TopicID:        "invalid-topic",
				},
			}
			Expect(k8sClient.Create(ctx, subscription)).Should(Succeed())
			subscriptionRef := types.NamespacedName{Namespace: subscription.Namespace, Name: subscription.Name}

			By("Checking if the status is Error")
			Eventually(func(g Gomega) {
				var subscription googlecloudpubsuboperatorv1.Subscription
				g.Expect(k8sClient.Get(ctx, subscriptionRef, &subscription)).Should(Succeed())
				g.Expect(subscription.Status.Phase).Should(Equal(googlecloudpubsuboperatorv1.SubscriptionStatusPhaseError))
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())
		})
	})
})
