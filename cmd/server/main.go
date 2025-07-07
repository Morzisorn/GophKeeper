package main

import (
	"gophkeeper/config"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server"
	"gophkeeper/internal/server/crypto"
	"gophkeeper/internal/server/repositories"
	cserv "gophkeeper/internal/server/services/crypto_service"
	userv "gophkeeper/internal/server/services/user_service"

	"go.uber.org/zap"
)

func main() {
	err := logger.Init()
	if err != nil {
		panic(err)
	}
	cnfg := config.GetServerConfig()
	err = crypto.LoadRSAKeyPair()
	if err != nil {
		logger.Log.Panic("failed to load RSA keys", zap.Error(err))
	}

	repo, err := repositories.NewStorage(cnfg)
	if err != nil {
		logger.Log.Panic("failed to create storage: %w", zap.Error(err))
	}
	us := userv.NewUserService(repo)
	cs := cserv.NewCryptoService(repo)

	server.CreateAndRun(us, cs, &repo)
}
