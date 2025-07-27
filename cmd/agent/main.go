package main

import (
	"fmt"
	"gophkeeper/config"
	"gophkeeper/internal/agent/client"
	"gophkeeper/internal/agent/services"
	"gophkeeper/internal/agent/ui"
	"gophkeeper/internal/logger"
	"os"
)

func main() {
	if err := runAgent(); err != nil {
		fmt.Printf("run agent error: %v\n", err)
		os.Exit(1)
	}
}

func runAgent() error {
	if err := logger.Init(); err != nil {
		return fmt.Errorf("init logger error: %w\n", err)
	}

	cnfg, err := config.GetAgentConfig()
	if err != nil {
		return fmt.Errorf("get agent config error: %w\n", err)
	}

	clnt, err := client.NewGRPCClient(cnfg)
	if err != nil {
		return fmt.Errorf("new grpc client error: %w\n", err)
	}

	cs, err := services.NewCryptoService(clnt)
	if err != nil {
		return fmt.Errorf("new crypto service error: %w\n", err)
	}

	us, err := services.NewUserService(clnt, cs)
	if err != nil {
		return fmt.Errorf("new user service error: %w\n", err)
	}

	is, err := services.NewItemService(clnt, cs)
	if err != nil {
		return fmt.Errorf("new item service error: %w\n", err)
	}

	if err = cs.SetPublicKey(); err != nil {
		return fmt.Errorf("set public key error: %w\n", err)
	}

	uiContr, err := ui.NewUIController(us, is)
	if err != nil {
		return fmt.Errorf("new ui controller error: %w\n", err)
	}
	if err := uiContr.Run(); err != nil {
		return fmt.Errorf("run ui controller error: %w\n", err)
	}
	return nil
}
