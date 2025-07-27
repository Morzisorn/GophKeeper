package crypto_service

import "gophkeeper/internal/server/repositories"

type CryptoService struct {
	repo repositories.Storage
}

func NewCryptoService(repo repositories.Storage) (*CryptoService, error) {
	return &CryptoService{repo: repo}, nil
}
