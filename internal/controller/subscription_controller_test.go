package controller

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//+kubebuilder:scaffold:imports
)

var _ = Describe("Subscription controller", func() {
	Context("When creating a Subscription resource", func() {
		const projectID = "subscription-project"
		It("Should create a Pub/Sub Subscription", func(ctx context.Context) {
			psClient, err := pubsub.NewClient(ctx, projectID,
				option.WithEndpoint(psServer.Addr),
				option.WithoutAuthentication(),
				option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
			)
			Expect(err).ShouldNot(HaveOccurred())

			topicID := "my-topic"
			_, err = psClient.CreateTopic(ctx, topicID)
			Expect(err).ShouldNot(HaveOccurred())

			By("Creating a Subscription")
			topic := &googlecloudpubsuboperatorv1.Subscription{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "googlecloudpubsuboperator.quipper.github.io/v1",
					Kind:       "Subscription",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: "default",
				},
				Spec: googlecloudpubsuboperatorv1.SubscriptionSpec{
					SubscriptionProjectID: projectID,
					SubscriptionID:        "my-subscription",
					TopicProjectID:        projectID,
					TopicID:               topicID,
				},
			}
			Expect(k8sClient.Create(ctx, topic)).Should(Succeed())

			By("Checking if the Subscription exists")
			Eventually(func(g Gomega) {
				subscriptionExists, err := psClient.Subscription("my-subscription").Exists(ctx)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(subscriptionExists).Should(BeTrue())
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())
		})
	})
})
