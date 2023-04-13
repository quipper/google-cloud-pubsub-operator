package controllers

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//+kubebuilder:scaffold:imports
)

var _ = Describe("Topic controller", func() {
	Context("When creating a Topic resource", func() {
		It("Should create a Pub/Sub Topic", func(ctx context.Context) {
			psClient, err := pubsub.NewClient(ctx, "my-project",
				option.WithEndpoint(psServer.Addr),
				option.WithoutAuthentication(),
				option.WithGRPCDialOption(grpc.WithInsecure()),
			)
			Expect(err).ShouldNot(HaveOccurred())

			By("Creating a Topic")
			topic := &googlecloudpubsuboperatorv1.Topic{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "googlecloudpubsuboperator.quipper.github.io/v1",
					Kind:       "Topic",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "example",
					Namespace: "default",
				},
				Spec: googlecloudpubsuboperatorv1.TopicSpec{
					ProjectID: "my-project",
					TopicID:   "my-topic",
				},
			}
			Expect(k8sClient.Create(ctx, topic)).Should(Succeed())

			By("Checking if the Topic exists")
			Eventually(func(g Gomega) {
				topicExists, err := psClient.Topic("my-topic").Exists(ctx)
				g.Expect(err).ShouldNot(HaveOccurred())
				g.Expect(topicExists).Should(BeTrue())
			}, 3*time.Second, 100*time.Millisecond).Should(Succeed())
		})
	})
})
