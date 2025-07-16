package services

import (
	"context"
	"gophkeeper/internal/agent/client"
	"gophkeeper/models"
)

type ItemService struct {
	Client client.Client
}

func NewItemService(client client.Client) *ItemService {
	return &ItemService{
		Client: client,
	}
}

func (ic *ItemService) AddItem(ctx context.Context, item *models.Item) error {
	return ic.Client.AddItem(ctx, item)
}

func (ic *ItemService) EditItem(ctx context.Context, item *models.Item) error {
	return ic.Client.EditItem(ctx, item)
}

func (ic *ItemService) DeleteItem(ctx context.Context, login, itemID string) error {
	return ic.Client.DeleteItem(ctx, login, itemID)
}

func (ic *ItemService) GetItems(ctx context.Context, login string, typ models.ItemType) ([]models.Item, error) {
	return ic.Client.GetItems(ctx, login, typ)
}

func (ic *ItemService) GetTypesCounts(ctx context.Context, login string) (map[string]int32, error) {
	return ic.Client.GetTypesCounts(ctx, login)
}
