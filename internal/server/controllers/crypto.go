package controllers

import (
	"context"
	"gophkeeper/config"
	pb "gophkeeper/internal/protos/crypto"
)

type CryptoController struct {
	pb.UnimplementedCryptoControllerServer
	cnfg config.ServerControllersConfig
}

func NewCryptoController(cnfg config.ServerControllersConfig) *CryptoController {
	return &CryptoController{cnfg: cnfg}
}

func (cc *CryptoController) GetPublicKeyPEM(ctx context.Context, in *pb.GetPublicKeyPEMRequest) (*pb.GetPublicKeyPEMResponse, error) {
	return &pb.GetPublicKeyPEMResponse{
		PublicKeyPem: string(cc.cnfg.GetPublicKeyPEM()),
	}, nil
}
