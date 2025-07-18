package controllers

import (
	"context"
	"gophkeeper/config"
	pb "gophkeeper/internal/protos/crypto"
)

type CryptoController struct {
	pb.UnimplementedCryptoControllerServer
	// service *crypto.CryptoService
}

func NewCryptoController() *CryptoController {
	return &CryptoController{}
}

func (cc *CryptoController) GetPublicKeyPEM(ctx context.Context, in *pb.GetPublicKeyPEMRequest) (*pb.GetPublicKeyPEMResponse, error) {
	return &pb.GetPublicKeyPEMResponse{
		PublicKeyPem: string(config.GetServerConfig().PublicKeyPEM),
	}, nil
}
