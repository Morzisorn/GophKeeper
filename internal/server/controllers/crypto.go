package controllers

import (
	"context"
	"fmt"
	"gophkeeper/config"
	pb "gophkeeper/internal/protos/crypto"
)

type CryptoController struct {
	pb.UnimplementedCryptoControllerServer
}

func NewCryptoController() *CryptoController {
	return &CryptoController{}
}

func (cc *CryptoController) GetPublicKeyPEM(ctx context.Context, in *pb.GetPublicKeyPEMRequest) (*pb.GetPublicKeyPEMResponse, error) {
	cnfg, err := config.GetServerConfig()
	if err != nil {
		return nil, fmt.Errorf("GetPublicKeyPEM: failed to get config: %w", err)
	}
	return &pb.GetPublicKeyPEMResponse{
		PublicKeyPem: string(cnfg.PublicKeyPEM),
	}, nil
}
