package item_service

import (
	"context"
	"fmt"
	"gophkeeper/internal/server/repositories"
	"gophkeeper/models"
)

type ItemService struct {
	repo repositories.Storage
}

func NewItemService(repo repositories.Storage) (*ItemService, error) {
	return &ItemService{repo: repo}, nil
}

func (is *ItemService) GetUserItems(ctx context.Context, typ models.ItemType, login string) ([]models.EncryptedItem, error) {
	var sl []models.EncryptedItem
	var err error
	if typ != models.ItemTypeUNSPECIFIED {
		sl, err = is.repo.GetUserItemsWithType(ctx, typ, login)
	} else {
		sl, err = is.repo.GetAllUserItems(ctx, login)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get %s from db for %s: %w", typ, login, err)
	}

	return sl, nil
}

func (is *ItemService) GetTypesCounts(ctx context.Context, login string) (map[models.ItemType]int32, error) {
	typesCount, err := is.repo.GetTypesCounts(ctx, login)
	if err != nil {
		return nil, err
	}

	if len(typesCount) != len(models.ItemTypes) {
		for _, t := range models.ItemTypes {
			_, ok := typesCount[t]
			if !ok {
				typesCount[t] = int32(0)
				continue
			}
		}
	}

	return typesCount, nil
}

func (is *ItemService) AddItem(ctx context.Context, item *models.EncryptedItem) error {
	return is.repo.AddItem(ctx, item)
}

func (is *ItemService) EditItem(ctx context.Context, item *models.EncryptedItem) error {
	return is.repo.EditItem(ctx, item)
}

func (is *ItemService) DeleteItem(ctx context.Context, login string, itemID [16]byte) error {
	return is.repo.DeleteItem(ctx, login, itemID)
}
