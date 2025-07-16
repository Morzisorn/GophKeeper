package database

import (
	"context"
	"encoding/json"
	"fmt"
	"gophkeeper/internal/errs"
	gen "gophkeeper/internal/server/repositories/database/generated"
	"gophkeeper/models"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemDatabase interface {
	GetAllUserItems(ctx context.Context, login string) ([]models.Item, error)
	GetUserItemsWithType(ctx context.Context, typ models.ItemType, login string) ([]models.Item, error)
	GetTypesCounts(ctx context.Context, login string) (map[models.ItemType]int32, error)
	AddItem(ctx context.Context, item *models.Item) error
	EditItem(ctx context.Context, item *models.Item) error
	DeleteItem(ctx context.Context, login string, itemID string) error
}

type ItemDB struct {
	q    *gen.Queries
	pool *pgxpool.Pool
}

func NewItemDB(q *gen.Queries, pool *pgxpool.Pool) ItemDatabase {
	return &ItemDB{
		q:    q,
		pool: pool,
	}
}

func (db *ItemDB) GetAllUserItems(ctx context.Context, login string) ([]models.Item, error) {
	dbItems, err := db.q.GetAllUserItems(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get all user items error: %w", err)
	}
	items := make([]models.Item, len(dbItems))
	for i, d := range dbItems {
		var meta models.Meta
		err := json.Unmarshal(d.Meta, &meta)
		if err != nil {
			return nil, fmt.Errorf("unmarshal meta info error: %w", err)
		}

		item := models.Item{
			ID:        d.ID.String(),
			Name:      d.Name,
			Meta:      meta,
			CreatedAt: d.CreatedAt.Time,
			UpdatedAt: d.UpdatedAt.Time,
		}

		switch models.ItemType(d.Type) {
		case models.ItemTypeCREDENTIALS:
			item.Data = &models.Credentials{
				Login:    d.Login.String,
				Password: d.Password.String,
			}
		case models.ItemTypeCARD:
			item.Data = &models.Card{
				Number:         d.Number.String,
				ExpiryDate:     d.ExpiryDate.String,
				SecurityCode:   d.SecurityCode.String,
				CardholderName: d.CardholderName.String,
			}
		case models.ItemTypeTEXT:
			item.Data = &models.Text{
				Content: d.TextContent.String,
			}
		case models.ItemTypeBINARY:
			item.Data = &models.Binary{
				Content: d.BinaryContent,
			}
		}
		items[i] = item
	}
	return items, nil
}

func (db *ItemDB) GetUserItemsWithType(ctx context.Context, typ models.ItemType, login string) ([]models.Item, error) {
	switch strings.ToUpper(typ.String()) {
	case string(gen.ItemTypeCREDENTIALS):
		return db.getCredentials(ctx, login)
	case string(gen.ItemTypeTEXT):
		return db.getTexts(ctx, login)
	case string(gen.ItemTypeBINARY):
		return db.getBinaries(ctx, login)
	case string(gen.ItemTypeCARD):
		return db.getCards(ctx, login)
	default:
		return nil, errs.ErrIncorrectItemType
	}
}

func (db *ItemDB) getCredentials(ctx context.Context, login string) ([]models.Item, error) {
	credsDB, err := db.q.GetCredentials(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get credentials from db error: %w", err)
	}

	creds := make([]models.Item, len(credsDB))
	for i, c := range credsDB {
		var meta models.Meta
		err := json.Unmarshal(c.Meta, &meta)
		if err != nil {
			return nil, fmt.Errorf("unmarshal meta info error: %w", err)
		}

		data := models.Credentials{
			Login:    c.Login.String,
			Password: c.Password.String,
		}

		creds[i] = models.Item{
			ID:        c.ID.String(),
			Name:      c.Name,
			Data:      data,
			Meta:      meta,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
		}
	}

	return creds, nil
}

func (db *ItemDB) getTexts(ctx context.Context, login string) ([]models.Item, error) {
	textsDB, err := db.q.GetTexts(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get texts from db error: %w", err)
	}

	texts := make([]models.Item, len(textsDB))
	for i, t := range textsDB {
		var meta models.Meta
		err := json.Unmarshal(t.Meta, &meta)
		if err != nil {
			return nil, fmt.Errorf("unmarshal meta info error: %w", err)
		}

		data := models.Text{Content: t.Content.String}

		texts[i] = models.Item{
			ID:        t.ID.String(),
			Name:      t.Name,
			Data:      data,
			Meta:      meta,
			CreatedAt: t.CreatedAt.Time,
			UpdatedAt: t.UpdatedAt.Time,
		}
	}

	return texts, nil
}

func (db *ItemDB) getBinaries(ctx context.Context, login string) ([]models.Item, error) {
	binsDB, err := db.q.GetBinaries(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get binaries from db error: %w", err)
	}

	bins := make([]models.Item, len(binsDB))
	for i, b := range binsDB {
		var meta models.Meta
		err := json.Unmarshal(b.Meta, &meta)
		if err != nil {
			return nil, fmt.Errorf("unmarshal meta info error: %w", err)
		}

		data := models.Binary{
			Content: b.Content,
		}

		bins[i] = models.Item{
			ID:        b.ID.String(),
			Name:      b.Name,
			Data:      data,
			Meta:      meta,
			CreatedAt: b.CreatedAt.Time,
			UpdatedAt: b.UpdatedAt.Time,
		}
	}

	return bins, nil
}

func (db *ItemDB) getCards(ctx context.Context, login string) ([]models.Item, error) {
	cardsDB, err := db.q.GetCards(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get cards from db error: %w", err)
	}

	cards := make([]models.Item, len(cardsDB))
	for i, c := range cardsDB {
		var meta models.Meta
		err := json.Unmarshal(c.Meta, &meta)
		if err != nil {
			return nil, fmt.Errorf("unmarshal meta info error: %w", err)
		}

		data := models.Card{
			Number:         c.Number.String,
			ExpiryDate:     c.ExpiryDate.String,
			SecurityCode:   c.SecurityCode.String,
			CardholderName: c.CardholderName.String,
		}

		cards[i] = models.Item{
			ID:        c.ID.String(),
			Name:      c.Name,
			Data:      data,
			Meta:      meta,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
		}
	}

	return cards, nil
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

func (db *ItemDB) AddItem(ctx context.Context, item *models.Item) error {
	err := withTransaction(ctx, db.pool, func(qtx *gen.Queries) error {
		meta, err := json.Marshal(item.Meta)
		if err != nil {
			return fmt.Errorf("marshal meta info error: %w", err)
		}
		itemID, err := qtx.AddItem(ctx, gen.AddItemParams{
			UserLogin: item.UserLogin,
			Name:      item.Name,
			Type:      gen.ItemType(item.Type),
			Meta:      meta,
		})
		if err != nil {
			return fmt.Errorf("add item error: %w", err)
		}

		switch item.Type {
		case models.ItemTypeCREDENTIALS:
			return db.addCredentials(ctx, qtx, itemID, item.Data.(*models.Credentials))
		case models.ItemTypeTEXT:
			return db.addText(ctx, itemID, item.Data.(*models.Text))
		case models.ItemTypeBINARY:
			return db.addBinary(ctx, itemID, item.Data.(*models.Binary))
		case models.ItemTypeCARD:
			return db.addCard(ctx, itemID, item.Data.(*models.Card))
		default:
			return errs.ErrIncorrectItemType
		}
	})
	if err != nil {
		return fmt.Errorf("add %s item error: %w", item.Type, err)
	}
	return nil
}

func (db *ItemDB) addCredentials(ctx context.Context, qtx *gen.Queries, itemID pgtype.UUID, cr *models.Credentials) error {
	return qtx.AddCredentials(ctx, gen.AddCredentialsParams{
		ItemID:   itemID,
		Login:    cr.Login,
		Password: cr.Password,
	})
}

func (db *ItemDB) addText(ctx context.Context, itemID pgtype.UUID, t *models.Text) error {
	return db.q.AddText(ctx, gen.AddTextParams{
		ItemID:  itemID,
		Content: t.Content,
	})
}

func (db *ItemDB) addBinary(ctx context.Context, itemID pgtype.UUID, b *models.Binary) error {
	return db.q.AddBinary(ctx, gen.AddBinaryParams{
		ItemID:  itemID,
		Content: b.Content,
	})
}

func (db *ItemDB) addCard(ctx context.Context, itemID pgtype.UUID, c *models.Card) error {
	return db.q.AddCard(ctx, gen.AddCardParams{
		ItemID:         itemID,
		Number:         c.Number,
		ExpiryDate:     c.ExpiryDate,
		SecurityCode:   c.SecurityCode,
		CardholderName: c.CardholderName,
	})
}

func (db *ItemDB) EditItem(ctx context.Context, item *models.Item) error {
	err := withTransaction(ctx, db.pool, func(qtx *gen.Queries) error {
		id, err := stringToPgUUID(item.ID)
		if err != nil {
			return fmt.Errorf("convert item id to UUID error: %w", err)
		}

		meta, err := json.Marshal(item.Meta)
		if err != nil {
			return fmt.Errorf("marshal meta info error: %w", err)
		}
		if err := qtx.EditItem(ctx, gen.EditItemParams{
			ID:   id,
			Name: item.Name,
			Meta: meta,
		}); err != nil {
			return fmt.Errorf("edit item error: %w", err)
		}

		switch item.Type {
		case models.ItemTypeCREDENTIALS:
			return db.editCredentials(ctx, id, item.Data.(*models.Credentials))
		case models.ItemTypeTEXT:
			return db.editText(ctx, id, item.Data.(*models.Text))
		case models.ItemTypeBINARY:
			return db.editBinary(ctx, id, item.Data.(*models.Binary))
		case models.ItemTypeCARD:
			return db.editCard(ctx, id, item.Data.(*models.Card))
		default:
			return errs.ErrIncorrectItemType
		}
	})
	if err != nil {
		return fmt.Errorf("edit %s item error: %w", item.Type, err)
	}
	return nil
}

func (db *ItemDB) editCredentials(ctx context.Context, itemID pgtype.UUID, cr *models.Credentials) error {
	return db.q.EditCredentials(ctx, gen.EditCredentialsParams{
		ItemID:   itemID,
		Login:    cr.Login,
		Password: cr.Password,
	})
}

func (db *ItemDB) editText(ctx context.Context, itemID pgtype.UUID, t *models.Text) error {
	return db.q.EditText(ctx, gen.EditTextParams{
		ItemID:  itemID,
		Content: t.Content,
	})
}

func (db *ItemDB) editBinary(ctx context.Context, itemID pgtype.UUID, b *models.Binary) error {
	return db.q.EditBinary(ctx, gen.EditBinaryParams{
		ItemID:  itemID,
		Content: b.Content,
	})
}

func (db *ItemDB) editCard(ctx context.Context, itemID pgtype.UUID, c *models.Card) error {
	return db.q.EditCard(ctx, gen.EditCardParams{
		ItemID:         itemID,
		Number:         c.Number,
		ExpiryDate:     c.ExpiryDate,
		SecurityCode:   c.SecurityCode,
		CardholderName: c.CardholderName,
	})
}

func (db *ItemDB) DeleteItem(ctx context.Context, login string, itemID string) error {
	id, err := stringToPgUUID(itemID)
	if err != nil {
		return fmt.Errorf("convert item id to UUID error: %w", err)
	}
	return db.q.DeleteItem(ctx, gen.DeleteItemParams{
		UserLogin: login,
		ID:        id,
	})
}

func withTransaction(ctx context.Context, db *pgxpool.Pool, fn func(q *gen.Queries) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	qtx := gen.New(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(qtx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func stringToPgUUID(s string) (pgtype.UUID, error) {
	var pgUUID pgtype.UUID
	err := pgUUID.Scan(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgUUID, nil
}
