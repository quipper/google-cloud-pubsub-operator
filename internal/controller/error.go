package controller

import (
	"errors"

	"google.golang.org/grpc/status"
)

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
