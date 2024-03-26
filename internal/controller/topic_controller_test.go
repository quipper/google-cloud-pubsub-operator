package controller

import (
	"context"
	"time"

	"cloud.google.com/go/iam"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
	"github.com/quipper/google-cloud-pubsub-operator/internal/pubsubtest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	//+kubebuilder:scaffold:imports
)

var _ = Describe("Topic controller", func() {
	Context("When creating a Topic resource", func() {
		It("Should create a Pub/Sub Topic", func(ctx context.Context) {
			const projectID = "topic-project-1"
			psClient, err := pubsubtest.NewClient(ctx, projectID, psServer)
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
				g.Expect(topic.Status.Message).Should(BeEmpty())
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())
		})

		It("Should create a Pub/Sub Topic with IAM Binding", func(ctx context.Context) {
			const projectID = "topic-project-1"
			const topicID = "my-topic-with-iam"

			psClient, err := pubsubtest.NewClient(ctx, projectID, psServer)
			Expect(err).ShouldNot(HaveOccurred())

			By("Creating a Topic with IAM Policy")
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
					TopicID:   topicID,
					Bindings: []googlecloudpubsuboperatorv1.IamBinding{
						{
							Role: "roles/pubsub.publisher",
							ServiceAccounts: []string{
								"my-app-1@topic-project-1.iam.gserviceaccount.com",
								"my-app-2@topic-project-1.iam.gserviceaccount.com",
							},
						},
						{
							Role: "roles/pubsub.viewer",
							ServiceAccounts: []string{
								"my-app-3@topic-project-1.iam.gserviceaccount.com",
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, topic)).Should(Succeed())
			topicRef := types.NamespacedName{Namespace: topic.Namespace, Name: topic.Name}

			By("Checking if the Topic exists")
			Eventually(func(g Gomega) {
				topic := psClient.Topic(topicID)
				topicExists, err := topic.Exists(ctx)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(topicExists).Should(BeTrue())

				policy, err := topic.IAM().Policy(ctx)
				g.Expect(err).ShouldNot(HaveOccurred())

				publishers := policy.Members(iam.RoleName("roles/pubsub.publisher"))
				viewers := policy.Members(iam.RoleName("roles/pubsub.viewer"))
				g.Expect(publishers).Should(ContainElement("my-app-1@topic-project-1.iam.gserviceaccount.com"))
				g.Expect(publishers).Should(ContainElement("my-app-2@topic-project-1.iam.gserviceaccount.com"))
				g.Expect(publishers).Should(HaveLen(2))
				g.Expect(viewers).Should(ContainElement("my-app-3@topic-project-1.iam.gserviceaccount.com"))
				g.Expect(viewers).Should(HaveLen(1))
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())

			By("Checking if the status is Active")
			Eventually(func(g Gomega) {
				var topic googlecloudpubsuboperatorv1.Topic
				g.Expect(k8sClient.Get(ctx, topicRef, &topic)).Should(Succeed())
				g.Expect(topic.Status.Phase).Should(Equal(googlecloudpubsuboperatorv1.TopicStatusPhaseActive))
				g.Expect(topic.Status.Message).Should(BeEmpty())
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())
		})

		It("Should update the status if error", func(ctx context.Context) {
			const projectID = "error-injected-project-topic-1"

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
					TopicID:   "this-should-be-failed",
				},
			}
			Expect(k8sClient.Create(ctx, topic)).Should(Succeed())
			topicRef := types.NamespacedName{Namespace: topic.Namespace, Name: topic.Name}

			By("Checking if the status is Error")
			Eventually(func(g Gomega) {
				var topic googlecloudpubsuboperatorv1.Topic
				g.Expect(k8sClient.Get(ctx, topicRef, &topic)).Should(Succeed())
				g.Expect(topic.Status.Phase).Should(Equal(googlecloudpubsuboperatorv1.TopicStatusPhaseError))
				g.Expect(topic.Status.Message).Should(HavePrefix("Pub/Sub error:"))
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())
		})
	})
})
