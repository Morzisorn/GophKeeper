package client

import (
	"context"
	"gophkeeper/config"
	"gophkeeper/internal/logger"
	pbcr "gophkeeper/internal/protos/crypto"
	pbus "gophkeeper/internal/protos/users"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	//User
	SignUpUser(ctx context.Context, login, password string) (token string, err error)
	SignInUser(ctx context.Context, login, password string) (token string, err error)
	SetJWTToken(token string)
	GetJWTToken() string

	//Crypto
	GetPublicKeyPEM(ctx context.Context) (string, error)
}

type GRPCClient struct {
	pbcr.CryptoControllerClient

	token   string
	conn    *grpc.ClientConn
	BaseURL string
	User    pbus.UserControllerClient
	Crypto  pbcr.CryptoControllerClient
}

// NewGRPCClient creates new pointer to GRPCClient based on config
func NewGRPCClient(c *config.Config) *GRPCClient {
	conn, err := grpc.NewClient(c.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Log.Fatal("create grpc client error: ", zap.Error(err))
	}
	return &GRPCClient{
		conn:    conn,
		BaseURL: c.Addr,
		User:    pbus.NewUserControllerClient(conn),
		Crypto:  pbcr.NewCryptoControllerClient(conn),
	}
}
