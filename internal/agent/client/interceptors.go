package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (g *GRPCClient) authInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if g.token != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+g.token)
	}

	return invoker(ctx, method, req, reply, cc, opts...)
}
