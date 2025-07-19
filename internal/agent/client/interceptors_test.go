package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGRPCClient_authInterceptor_WithToken(t *testing.T) {
	client := &GRPCClient{
		token: "test-token",
	}

	var capturedCtx context.Context
	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		capturedCtx = ctx
		return nil
	}

	ctx := context.Background()
	err := client.authInterceptor(ctx, "test-method", nil, nil, nil, mockInvoker)

	assert.NoError(t, err)
	
	// Проверяем, что токен добавлен в метаданные
	md, ok := metadata.FromOutgoingContext(capturedCtx)
	assert.True(t, ok)
	
	authValues := md.Get("authorization")
	assert.Len(t, authValues, 1)
	assert.Equal(t, "Bearer test-token", authValues[0])
}

func TestGRPCClient_authInterceptor_WithoutToken(t *testing.T) {
	client := &GRPCClient{
		token: "",
	}

	var capturedCtx context.Context
	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		capturedCtx = ctx
		return nil
	}

	ctx := context.Background()
	err := client.authInterceptor(ctx, "test-method", nil, nil, nil, mockInvoker)

	assert.NoError(t, err)
	
	// Проверяем, что метаданные не содержат токен
	md, ok := metadata.FromOutgoingContext(capturedCtx)
	if ok {
		authValues := md.Get("authorization")
		assert.Empty(t, authValues)
	}
}