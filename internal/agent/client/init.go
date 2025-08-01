package client

import (
	"context"
	"fmt"
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
	SetJWTToken(token string) error
	GetJWTToken() (string, error)

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
	token string
	conn  *grpc.ClientConn
	cnfg  config.AgentClientConfig

	User   pbus.UserControllerClient
	Crypto pbcr.CryptoControllerClient
	Item   pbit.ItemsControllerClient
}

func NewGRPCClient(cnfg config.AgentClientConfig) (*GRPCClient, error) {
	client := &GRPCClient{
		cnfg: cnfg,
	}

	conn, err := grpc.NewClient(cnfg.GetAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(client.authInterceptor),
	)
	if err != nil {
		logger.Log.Fatal("create grpc client error: ", zap.Error(err))
	}

	client.conn = conn
	client.User, err = pbus.NewUserControllerClient(conn)
	if err != nil {
		return nil, fmt.Errorf("create grpc client error: %w", err)
	}
	client.Crypto, err = pbcr.NewCryptoControllerClient(conn)
	if err != nil {
		return nil, fmt.Errorf("create grpc client error: %w", err)
	}
	client.Item, err = pbit.NewItemsControllerClient(conn)
	if err != nil {
		return nil, fmt.Errorf("create grpc client error: %w", err)
	}

	return client, nil

}
