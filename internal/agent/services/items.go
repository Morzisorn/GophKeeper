package services

import (
	"context"
	"fmt"
	"gophkeeper/internal/agent/client"
	"gophkeeper/models"
)

type ItemService struct {
	Client client.Client
	Crypto *CryptoService
}

func NewItemService(client client.Client, cs *CryptoService) (*ItemService, error) {
	return &ItemService{
		Client: client,
		Crypto: cs,
	}, nil
}

func (ic *ItemService) AddItem(ctx context.Context, item *models.Item) error {
	encItem, err := ic.Crypto.EncryptItem(item)
	if err != nil {
		return fmt.Errorf("encrypt item error: %w", err)
	}
	return ic.Client.AddItem(ctx, encItem)
}

func (ic *ItemService) EditItem(ctx context.Context, item *models.Item) error {
	encItem, err := ic.Crypto.EncryptItem(item)
	if err != nil {
		return fmt.Errorf("encrypt item error: %w", err)
	}
	return ic.Client.EditItem(ctx, encItem)
}

func (ic *ItemService) DeleteItem(ctx context.Context, login string, itemID [16]byte) error {
	return ic.Client.DeleteItem(ctx, login, itemID)
}

func (ic *ItemService) GetItems(ctx context.Context, login string, typ models.ItemType) ([]models.EncryptedItem, error) {
	return ic.Client.GetItems(ctx, login, typ)
}

func (is *ItemService) GetTypesCounts(ctx context.Context, login string) (map[string]int32, error) {
	return is.Client.GetTypesCounts(ctx, login)
}

func (is *ItemService) DecryptItem(encItem *models.EncryptedItem) (*models.Item, error) {
	return is.Crypto.DecryptItem(encItem)
}
