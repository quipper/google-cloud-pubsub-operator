package pubsubtest

import (
	"strings"

	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"cloud.google.com/go/pubsub/pstest"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateTopicErrorInjectionReactor() pstest.ServerReactorOption {
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

func CreateSubscriptionErrorInjectionReactor() pstest.ServerReactorOption {
	return pstest.ServerReactorOption{
		FuncName: "CreateSubscription",
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
