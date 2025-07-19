package server

import (
	"gophkeeper/config"
	"gophkeeper/internal/server/repositories/database"
	"gophkeeper/internal/server/services/crypto_service"
	"gophkeeper/internal/server/services/item_service"
	"gophkeeper/internal/server/services/user_service"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateGRPCServer(t *testing.T) {
	repo := &database.PGDB{}
	us := user_service.NewUserService(repo)
	cs := crypto_service.NewCryptoService(repo)
	is := item_service.NewItemService(repo)

	config.GetServerConfig().Addr = "127.0.0.1:12345"
	server := createGRPCServer(us, cs, is)
	require.NotNil(t, server)
}
