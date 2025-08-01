package server

import (
	"gophkeeper/config"
	"gophkeeper/internal/server/repositories"
	"gophkeeper/internal/server/repositories/database"
	"gophkeeper/internal/server/services/crypto_service"
	"gophkeeper/internal/server/services/item_service"
	"gophkeeper/internal/server/services/user_service"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateGRPCServer(t *testing.T) {
	cnfg, err := config.NewServerConfig()
	require.NoError(t, err)

	// Create a mock storage
	var repo repositories.Storage = &database.PGDB{}

	us, err := user_service.NewUserService(cnfg, repo)
	require.NoError(t, err)
	cs, err := crypto_service.NewCryptoService(cnfg)
	require.NoError(t, err)
	is, err := item_service.NewItemService(repo)
	require.NoError(t, err)

	server, err := createGRPCServer(cnfg, us, cs, is)
	require.NoError(t, err)
	require.NotNil(t, server)
}
