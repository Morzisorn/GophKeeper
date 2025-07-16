package client

import (
	"context"
	"gophkeeper/config"
	"gophkeeper/internal/logger"
	pbcr "gophkeeper/internal/protos/crypto"
	pbit "gophkeeper/internal/protos/items"
	pbus "gophkeeper/internal/protos/users"
	"gophkeeper/models"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	//User
	SignUpUser(ctx context.Context, user *models.User) (token string, salt string, err error)
	SignInUser(ctx context.Context, user *models.User) (token string, salt string, err error)
	SetJWTToken(token string)
	GetJWTToken() string

	//Crypto
	GetPublicKeyPEM(ctx context.Context) (string, error)

	//Items
	AddItem(ctx context.Context, item *models.Item) error
	EditItem(ctx context.Context, item *models.Item) error
	DeleteItem(ctx context.Context, login, itemID string) error 
	GetItems(ctx context.Context, login string, typ models.ItemType) ([]models.Item, error)
	GetTypesCounts(ctx context.Context, login string) (map[string]int32, error) 
}

type GRPCClient struct {
	//pbcr.CryptoControllerClient

	token   string
	conn    *grpc.ClientConn
	BaseURL string
	User    pbus.UserControllerClient
	Crypto  pbcr.CryptoControllerClient
	Item    pbit.ItemsControllerClient
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
		Item:    pbit.NewItemsControllerClient(conn),
	}
}
