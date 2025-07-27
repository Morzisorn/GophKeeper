package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	gen "gophkeeper/internal/server/repositories/database/generated"
	"gophkeeper/models"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type ItemDatabase interface {
	GetAllUserItems(ctx context.Context, login string) ([]models.EncryptedItem, error)
	GetUserItemsWithType(ctx context.Context, typ models.ItemType, login string) ([]models.EncryptedItem, error)
	GetTypesCounts(ctx context.Context, login string) (map[models.ItemType]int32, error)
	AddItem(ctx context.Context, item *models.EncryptedItem) error
	EditItem(ctx context.Context, item *models.EncryptedItem) error
	DeleteItem(ctx context.Context, login string, itemID [16]byte) error
}

type PoolInterface interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

var _ PoolInterface = (*pgxpool.Pool)(nil)

type ItemDB struct {
	q    *gen.Queries
	pool PoolInterface
}

var _ ItemDatabase = (*ItemDB)(nil)

func NewItemDB(q *gen.Queries, pool PoolInterface) (ItemDatabase, error) {
	if pool == nil || q == nil {
		return nil, errors.New("create user database error: pool or quaries is nil")
	}
	return &ItemDB{
		q:    q,
		pool: pool,
	}, nil
}

func (db *ItemDB) GetAllUserItems(ctx context.Context, login string) ([]models.EncryptedItem, error) {
	dbItems, err := db.q.GetAllUserItems(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get all user items error: %w", err)
	}
	items := make([]models.EncryptedItem, len(dbItems))
	for i, d := range dbItems {
		var meta models.Meta
		err := json.Unmarshal(d.Meta, &meta)
		if err != nil {
			return nil, fmt.Errorf("unmarshal meta info error: %w", err)
		}

		encData := models.EncryptedData{
			EncryptedContent: d.EncryptedDataContent,
			Nonce:            d.EncryptedDataNonce,
		}

		items[i] = models.EncryptedItem{
			ID:            d.ID.Bytes,
			UserLogin:     login,
			Name:          d.Name,
			Type:          models.ItemType(d.Type),
			EncryptedData: encData,
			Meta:          meta,
			CreatedAt:     d.CreatedAt.Time,
			UpdatedAt:     d.UpdatedAt.Time,
		}
	}
	return items, nil
}

func (db *ItemDB) GetUserItemsWithType(ctx context.Context, typ models.ItemType, login string) ([]models.EncryptedItem, error) {
	dbItems, err := db.q.GetUserItemsWithType(ctx, gen.GetUserItemsWithTypeParams{
		UserLogin: login,
		Type:      itemTypeModelsToPg(typ),
	})
	if err != nil {
		return nil, fmt.Errorf("get user items with type %s error: %w", typ.String(), err)
	}

	items := make([]models.EncryptedItem, len(dbItems))
	for i, d := range dbItems {
		var meta models.Meta
		err := json.Unmarshal(d.Meta, &meta)
		if err != nil {
			return nil, fmt.Errorf("unmarshal meta info error: %w", err)
		}

		encData := models.EncryptedData{
			EncryptedContent: d.EncryptedDataContent,
			Nonce:            d.EncryptedDataNonce,
		}

		items[i] = models.EncryptedItem{
			ID:            d.ID.Bytes,
			UserLogin:     login,
			Name:          d.Name,
			Type:          models.ItemType(d.Type),
			EncryptedData: encData,
			Meta:          meta,
			CreatedAt:     d.CreatedAt.Time,
			UpdatedAt:     d.UpdatedAt.Time,
		}
	}
	return items, nil
}

func itemTypeModelsToPg(typ models.ItemType) gen.ItemType {
	switch strings.ToUpper(typ.String()) {
	case string(gen.ItemTypeCREDENTIALS):
		return gen.ItemTypeCREDENTIALS
	case string(gen.ItemTypeTEXT):
		return gen.ItemTypeTEXT
	case string(gen.ItemTypeBINARY):
		return gen.ItemTypeBINARY
	case string(gen.ItemTypeCARD):
		return gen.ItemTypeCARD
	default:
		return "unknown"
	}
}

func (db *ItemDB) GetTypesCounts(ctx context.Context, login string) (map[models.ItemType]int32, error) {
	dbCounts, err := db.q.GetTypesCounts(ctx, login)
	res := make(map[models.ItemType]int32, len(dbCounts))
	if err != nil {
		return nil, fmt.Errorf("failed to get counts of each item type: %w", err)
	}

	for _, t := range dbCounts {
		res[models.ItemType(t.Type)] = int32(t.Count)
	}
	return res, nil
}

func (db *ItemDB) AddItem(ctx context.Context, item *models.EncryptedItem) error {
	meta, err := json.Marshal(item.Meta)
	if err != nil {
		return fmt.Errorf("marshal meta info error: %w", err)
	}
	if _, err = db.q.AddItem(ctx, gen.AddItemParams{
		UserLogin:            item.UserLogin,
		Name:                 item.Name,
		Type:                 gen.ItemType(item.Type),
		EncryptedDataContent: item.EncryptedData.EncryptedContent,
		EncryptedDataNonce:   item.EncryptedData.Nonce,
		Meta:                 meta,
	}); err != nil {
		return fmt.Errorf("add item error: %w", err)
	}
	return nil
}

func (db *ItemDB) EditItem(ctx context.Context, item *models.EncryptedItem) error {
	meta, err := json.Marshal(item.Meta)
	if err != nil {
		return fmt.Errorf("marshal meta info error: %w", err)
	}
	if err := db.q.EditItem(ctx, gen.EditItemParams{
		ID:                   pgtype.UUID{Bytes: item.ID, Valid: true},
		Name:                 item.Name,
		EncryptedDataContent: item.EncryptedData.EncryptedContent,
		EncryptedDataNonce:   item.EncryptedData.Nonce,
		Meta:                 meta,
	}); err != nil {
		return fmt.Errorf("edit item error: %w", err)
	}

	return nil
}

func (db *ItemDB) DeleteItem(ctx context.Context, login string, itemID [16]byte) error {
	return db.q.DeleteItem(ctx, gen.DeleteItemParams{
		UserLogin: login,
		ID:        pgtype.UUID{Bytes: itemID, Valid: true},
	})
}
