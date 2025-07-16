package client

import (
	"context"
	"fmt"
	"gophkeeper/models"

	pbit "gophkeeper/internal/protos/items"
)

func (g *GRPCClient) AddItem(ctx context.Context, item *models.Item) error {
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

func (g *GRPCClient) EditItem(ctx context.Context, item *models.Item) error {
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

func (g *GRPCClient) DeleteItem(ctx context.Context, login, itemID string) error {
	resp, err := g.Item.DeleteItem(ctx, &pbit.DeleteItemRequest{UserLogin: login, ItemId: itemID})
	if err != nil || !resp.Success {
		return fmt.Errorf("delete item server error: %w", err)
	}

	return nil
}

func (g *GRPCClient) GetItems(ctx context.Context, login string, typ models.ItemType) ([]models.Item, error) {
	req := pbit.GetUserItemsRequest{UserLogin: login}
	if typ != "" {
		req.Type = typ.ToPb()
	}
	resp, err := g.Item.GetUserItems(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("get user items server error: %w", err)
	}

	items := make([]models.Item, len(resp.Items))

	for i, pbItem := range resp.Items {
		item, err := models.ItemPbToModels(pbItem)
		if err != nil {
			return nil, fmt.Errorf("convert pb item to models error: %w", err)
		}
		items[i] = *item
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

