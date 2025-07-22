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
	AddItem(ctx context.Context, item *models.EncryptedItem) error
	EditItem(ctx context.Context, item *models.EncryptedItem) error
	DeleteItem(ctx context.Context, login string, itemID [16]byte) error
	GetItems(ctx context.Context, login string, typ models.ItemType) ([]models.EncryptedItem, error)
	GetTypesCounts(ctx context.Context, login string) (map[string]int32, error)
}

var _ Client = (*GRPCClient)(nil)

type GRPCClient struct {
	token   string
	conn    *grpc.ClientConn
	BaseURL string

	User   pbus.UserControllerClient
	Crypto pbcr.CryptoControllerClient
	Item   pbit.ItemsControllerClient
}

func NewGRPCClient(c *config.Config) *GRPCClient {
	client := &GRPCClient{
		BaseURL: c.Addr,
	}

	conn, err := grpc.NewClient(c.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(client.authInterceptor),
	)
	if err != nil {
		logger.Log.Fatal("create grpc client error: ", zap.Error(err))
	}

	client.conn = conn
	client.User = pbus.NewUserControllerClient(conn)
	client.Crypto = pbcr.NewCryptoControllerClient(conn)
	client.Item = pbit.NewItemsControllerClient(conn)

	return client

}

