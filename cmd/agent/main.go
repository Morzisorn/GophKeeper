package main

import (
	"gophkeeper/config"
	"gophkeeper/internal/agent/client"
	"gophkeeper/internal/agent/services"
	"gophkeeper/internal/agent/ui"
	"gophkeeper/internal/logger"

	"go.uber.org/zap"
)


func main() {
	err := logger.Init()
	if err != nil {
		panic(err)
	}
	cnfg := config.GetAgentConfig()
	client := client.NewGRPCClient(cnfg)
	cs := services.NewCryptoService(client)
	us := services.NewUserService(client, cs)
	is := services.NewItemService(client, cs)
	err = cs.SetPublicKey()
	if err != nil {
		logger.Log.Fatal("Set public key error: ", zap.Error(err))
	}
	uiContr := ui.NewUIController(us, is)
	uiContr.Run()
}