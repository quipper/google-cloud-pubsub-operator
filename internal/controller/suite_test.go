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
	"path/filepath"
	"strings"
	"testing"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"cloud.google.com/go/pubsub/pstest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	googlecloudpubsuboperatorv1 "github.com/quipper/google-cloud-pubsub-operator/api/v1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var psServer *pstest.Server

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
	}

	ctx, cancel := context.WithCancel(context.TODO())
	DeferCleanup(func() {
		cancel()
		By("tearing down the test environment")
		Expect(testEnv.Stop()).Should(Succeed())
	})

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = googlecloudpubsuboperatorv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	psServer = pstest.NewServer(
		createTopicErrorInjectionReactor(),
	)
	DeferCleanup(func() {
		Expect(psServer.Close()).Should(Succeed())
	})

	newClient := func(ctx context.Context, projectID string, opts ...option.ClientOption) (c *pubsub.Client, err error) {
		return pubsub.NewClient(ctx, projectID,
			append(opts,
				option.WithEndpoint(psServer.Addr),
				option.WithoutAuthentication(),
				option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
			)...,
		)
	}

	err = (&TopicReconciler{
		Client:    k8sManager.GetClient(),
		Scheme:    k8sManager.GetScheme(),
		NewClient: newClient,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&SubscriptionReconciler{
		Client:    k8sManager.GetClient(),
		Scheme:    k8sManager.GetScheme(),
		NewClient: newClient,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})

func createTopicErrorInjectionReactor() pstest.ServerReactorOption {
	return pstest.ServerReactorOption{
		FuncName: "CreateTopic",
		Reactor: errorInjectionReactorFunc(func(req interface{}) (bool, interface{}, error) {
			topic, ok := req.(*pubsubpb.Topic)
			if ok {
				if strings.HasPrefix(topic.Name, "projects/error-injected-") {
					return true, nil, status.Errorf(codes.InvalidArgument, "error injected")
				}
			}
			return false, nil, nil
		}),
	}
}

type errorInjectionReactorFunc func(req interface{}) (bool, interface{}, error)

func (r errorInjectionReactorFunc) React(req interface{}) (bool, interface{}, error) {
	return r(req)
}
