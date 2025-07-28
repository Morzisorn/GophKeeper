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
	us, err := user_service.NewUserService(repo)
	require.NoError(t, err)
	cs, err := crypto_service.NewCryptoService(repo)
	require.NoError(t, err)
	is, err := item_service.NewItemService(repo)
	require.NoError(t, err)

	cnfg, err := config.GetServerConfig()
	require.NoError(t, err)
	cnfg.Addr = "127.0.0.1:12345"
	server, err := createGRPCServer(us, cs, is)
	require.NoError(t, err)
	require.NotNil(t, server)
}
