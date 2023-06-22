package controller

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	//+kubebuilder:scaffold:imports
)

var _ = Describe("Topic controller", func() {
	Context("When creating a Topic resource", func() {
		const projectID = "topic-project"
		It("Should create a Pub/Sub Topic", func(ctx context.Context) {
			psClient, err := pubsub.NewClient(ctx, projectID,
				option.WithEndpoint(psServer.Addr),
				option.WithoutAuthentication(),
				option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
			)
			Expect(err).ShouldNot(HaveOccurred())

			By("Creating a Topic")
			topic := &googlecloudpubsuboperatorv1.Topic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "googlecloudpubsuboperator.quipper.github.io/v1",
					Kind:       "Topic",
				},
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "example-",
					Namespace:    "default",
				},
				Spec: googlecloudpubsuboperatorv1.TopicSpec{
					ProjectID: projectID,
					TopicID:   "my-topic",
				},
			}
			Expect(k8sClient.Create(ctx, topic)).Should(Succeed())
			topicRef := types.NamespacedName{Namespace: topic.Namespace, Name: topic.Name}

			By("Checking if the Topic exists")
			Eventually(func(g Gomega) {
				topicExists, err := psClient.Topic("my-topic").Exists(ctx)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(topicExists).Should(BeTrue())
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())

			By("Checking if the status is Active")
			Eventually(func(g Gomega) {
				var topic googlecloudpubsuboperatorv1.Topic
				g.Expect(k8sClient.Get(ctx, topicRef, &topic)).Should(Succeed())
				g.Expect(topic.Status.Phase).Should(Equal(googlecloudpubsuboperatorv1.TopicStatusPhaseActive))
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())
		})

		It("Should update the status if error", func(ctx context.Context) {
			By("Creating a Topic")
			topic := &googlecloudpubsuboperatorv1.Topic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "googlecloudpubsuboperator.quipper.github.io/v1",
					Kind:       "Topic",
				},
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "example-",
					Namespace:    "default",
				},
				Spec: googlecloudpubsuboperatorv1.TopicSpec{
					ProjectID: projectID,
					// FIXME: fake server does not reject "goog" topic!
					// https://cloud.google.com/pubsub/docs/admin#resource_names
					TopicID: "goog-this-is-invalid",
				},
			}
			Expect(k8sClient.Create(ctx, topic)).Should(Succeed())
			topicRef := types.NamespacedName{Namespace: topic.Namespace, Name: topic.Name}

			By("Checking if the status is Error")
			Eventually(func(g Gomega) {
				var topic googlecloudpubsuboperatorv1.Topic
				g.Expect(k8sClient.Get(ctx, topicRef, &topic)).Should(Succeed())
				g.Expect(topic.Status.Phase).Should(Equal(googlecloudpubsuboperatorv1.TopicStatusPhaseError))
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())
		})
	})
})
