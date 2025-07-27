package main

import (
	"fmt"
	"gophkeeper/config"
	"gophkeeper/internal/logger"
	"gophkeeper/internal/server"
	"gophkeeper/internal/server/crypto"
	"gophkeeper/internal/server/repositories"
	cserv "gophkeeper/internal/server/services/crypto_service"
	iserv "gophkeeper/internal/server/services/item_service"
	userv "gophkeeper/internal/server/services/user_service"
	"os"
)

func main() {
	if err := runServer(); err != nil {
		fmt.Printf("run server error: %v\n", err)
		os.Exit(1)
	}
}

func runServer() error {
	if err := logger.Init(); err != nil {
		return fmt.Errorf("init logger error: %w\n", err)
	}
	cnfg, err := config.GetServerConfig()
	if err != nil {
		return fmt.Errorf("get agent config error: %w\n", err)
	}
	if err := crypto.LoadRSAKeyPair(); err != nil {
		return fmt.Errorf("failed to load RSA keys: %w\n", err)
	}

	repo, err := repositories.NewStorage(cnfg)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w\n", err)
	}

	us, err := userv.NewUserService(repo)
	if err != nil {
		return fmt.Errorf("failed to create user service: %w\n", err)
	}
	cs, err := cserv.NewCryptoService(repo)
	if err != nil {
		return fmt.Errorf("failed to create crypto service: %w\n", err)
	}
	ic, err := iserv.NewItemService(repo)
	if err != nil {
		return fmt.Errorf("failed to create item service: %w\n", err)
	}

	if err := server.CreateAndRun(us, cs, ic); err != nil {
		return fmt.Errorf("create server error: %w\n", err)
	}

	return nil
}
