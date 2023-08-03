package controller

import (
	"context"
	"errors"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type newPubSubClientFunc func(ctx context.Context, projectID string, opts ...option.ClientOption) (*pubsub.Client, error)

func isPubSubNotFoundError(err error) bool {
	if gs, ok := gRPCStatusFromError(err); ok {
		return gs.Code() == codes.NotFound
	}
	return false
}

func isPubSubAlreadyExistsError(err error) bool {
	if gs, ok := gRPCStatusFromError(err); ok {
		return gs.Code() == codes.AlreadyExists
	}
	return false
}

func gRPCStatusFromError(err error) (*status.Status, bool) {
	type gRPCError interface {
		GRPCStatus() *status.Status
	}

	var se gRPCError
	if errors.As(err, &se) {
		return se.GRPCStatus(), true
	}

	return nil, false
}
