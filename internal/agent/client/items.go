package client

import (
	"context"
	"fmt"
	"gophkeeper/models"

	pbit "gophkeeper/internal/protos/items"
)

func (g *GRPCClient) AddItem(ctx context.Context, item *models.EncryptedItem) error {
	pbItem, err := item.ToPb()
	if err != nil {
		return fmt.Errorf("convert model item to pb error: %w", err)
	}

	resp, err := g.Item.AddItem(ctx, &pbit.AddItemRequest{Item: pbItem})
	if err != nil || !resp.Success {
		return fmt.Errorf("add item server error: %w", err)
	}

	return nil
}

func (g *GRPCClient) EditItem(ctx context.Context, item *models.EncryptedItem) error {
	pbItem, err := item.ToPb()
	if err != nil {
		return fmt.Errorf("convert model item to pb error: %w", err)
	}

	resp, err := g.Item.EditItem(ctx, &pbit.EditItemRequest{Item: pbItem})
	if err != nil || !resp.Success {
		return fmt.Errorf("edit item server error: %w", err)
	}

	return nil
}

func (g *GRPCClient) DeleteItem(ctx context.Context, login string, itemID [16]byte) error {
	resp, err := g.Item.DeleteItem(ctx, &pbit.DeleteItemRequest{UserLogin: login, ItemId: itemID[:]})
	if err != nil || !resp.Success {
		return fmt.Errorf("delete item server error: %w", err)
	}

	return nil
}

func (g *GRPCClient) GetItems(ctx context.Context, login string, typ models.ItemType) ([]models.EncryptedItem, error) {
	req := pbit.GetUserItemsRequest{
		UserLogin: login,
		Type:      typ.ToPb(),
	}

	resp, err := g.Item.GetUserItems(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("get user items server error: %w", err)
	}

	items := make([]models.EncryptedItem, len(resp.Items))

	for i, pbItem := range resp.Items {
		items[i] = *models.EncryptedItemPbToModels(pbItem)
	}
	return items, nil
}

func (g *GRPCClient) GetTypesCounts(ctx context.Context, login string) (map[string]int32, error) {
	resp, err := g.Item.TypesCounts(ctx, &pbit.TypesCountsRequest{UserLogin: login})
	if err != nil {
		return nil, fmt.Errorf("get item type counters error: %w", err)
	}
	return resp.GetTypes(), nil
}
