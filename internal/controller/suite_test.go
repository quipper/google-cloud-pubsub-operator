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
	"reflect"
	"testing"
	"time"
	"unsafe"

	"cloud.google.com/go/iam/apiv1/iampb"
	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/quipper/google-cloud-pubsub-operator/internal/pubsubtest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	clocktesting "k8s.io/utils/clock/testing"
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

func TestControllers(t *testing.T) {
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
		pubsubtest.CreateTopicErrorInjectionReactor(),
	)

	// try to get access to `Gsrv` thru `*testutil.Server`, that's on first field of `pstest.Server`
	rs := reflect.ValueOf(psServer).Elem()

	// *testutil.Server
	tsrv := rs.Field(0)
	tsrv = reflect.NewAt(tsrv.Type(), unsafe.Pointer(tsrv.UnsafeAddr())).Elem()

	// *grpc.Server is found from the exported field `Gsrv`
	gsrv := reflect.Indirect(tsrv).FieldByName("Gsrv").Interface().(*grpc.Server)

	// trying to register fake iam policy server
	// this fails randomly because Registering is not possible after `grsv.Serve()` is called.
	iampb.RegisterIAMPolicyServer(gsrv, pubsubtest.CreateFakeIamPolicyServer())

	DeferCleanup(func() {
		Expect(psServer.Close()).Should(Succeed())
	})

	newClient := func(ctx context.Context, projectID string, opts ...option.ClientOption) (c *pubsub.Client, err error) {
		return pubsubtest.NewClient(ctx, projectID, psServer)
	}

	err = (&TopicReconciler{
		Client:    k8sManager.GetClient(),
		Scheme:    k8sManager.GetScheme(),
		NewClient: newClient,
		Recorder:  k8sManager.GetEventRecorderFor("topic-controller"),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	err = (&SubscriptionReconciler{
		Client:    k8sManager.GetClient(),
		Scheme:    k8sManager.GetScheme(),
		NewClient: newClient,
		Recorder:  k8sManager.GetEventRecorderFor("subscription-controller"),
		// TODO: how to change the time in test code?
		Clock: clocktesting.NewFakePassiveClock(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		defer GinkgoRecover()
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()
})
