package pubsubtest

import (
	"context"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewClient(ctx context.Context, projectID string, server *pstest.Server) (*pubsub.Client, error) {
	return pubsub.NewClient(ctx, projectID,
		option.WithEndpoint(server.Addr),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
}
