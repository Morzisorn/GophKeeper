package client

import (
	"context"
	"errors"
	"fmt"
	pb "gophkeeper/internal/protos/crypto"
)

func (g *GRPCClient) GetPublicKeyPEM(ctx context.Context) (string, error) {
	resp, err := g.Crypto.GetPublicKeyPEM(ctx, &pb.GetPublicKeyPEMRequest{})
	if err != nil {
		return "", fmt.Errorf("get public key pem error: %w", err)
	}

	if resp.PublicKeyPem == "" {
		return "", errors.New("pem is empty")
	}

	return resp.PublicKeyPem, nil
}
