package pubsubtest

import (
	"context"

	"cloud.google.com/go/iam/apiv1/iampb"
)

type FakeIamPolicyServer struct {
	iampb.UnimplementedIAMPolicyServer

	policyData map[string]*iampb.Policy
}

func CreateFakeIamPolicyServer() *FakeIamPolicyServer {
	return &FakeIamPolicyServer{
		policyData: make(map[string]*iampb.Policy),
	}
}

func (fake *FakeIamPolicyServer) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest) (*iampb.Policy, error) {
	if policy, exists := fake.policyData[req.Resource]; exists {
		return policy, nil
	} else {
		return &iampb.Policy{}, nil
	}
}

func (fake *FakeIamPolicyServer) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest) (*iampb.Policy, error) {
	fake.policyData[req.Resource] = req.Policy
	return req.Policy, nil
}
