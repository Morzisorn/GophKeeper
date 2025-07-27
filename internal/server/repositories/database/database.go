package database

import (
	"context"
	"fmt"
	"gophkeeper/config"
	"gophkeeper/models"
	"os"
	"path/filepath"

	gen "gophkeeper/internal/server/repositories/database/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	UserDatabase
	ItemDatabase
}

type PGDB struct {
	users UserDatabase
	items ItemDatabase
}

var _ Database = (*PGDB)(nil)

func NewPGDB(cfg *config.Config) (Database, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DBConnStr)
	if err != nil {
		return nil, fmt.Errorf("create new db error: %v", err)
	}

	if err := createTables(pool); err != nil {
		return nil, fmt.Errorf("create db tables error: %v", err)
	}

	q := gen.New(pool)

	userDB, err := NewUserDB(q, pool)
	if err != nil {
		return nil, fmt.Errorf("create user db error: %v", err)
	}
	itemDB, err := NewItemDB(q, pool)
	if err != nil {
		return nil, fmt.Errorf("create item db error: %v", err)
	}
	return &PGDB{
		users: userDB,
		items: itemDB,
	}, nil
}

func createTables(db *pgxpool.Pool) error {
	rootDir, err := config.GetProjectRoot()
	if err != nil {
		return err
	}
	dirpath := filepath.Join(rootDir, "internal", "server", "repositories", "database", "schema")

	var typeExists bool
	err = db.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'item_type')").Scan(&typeExists)
	if err != nil {
		return fmt.Errorf("failed to check if type exists: %w", err)
	}
	if !typeExists {
		filepathTypes := filepath.Join(dirpath, "001_types.sql")
		err = createTable(db, filepathTypes)
		if err != nil {
			return err
		}
	}

	filepathTables := filepath.Join(dirpath, "002_tables.sql")
	err = createTable(db, filepathTables)
	if err != nil {
		return err
	}

	return nil
}

func createTable(db *pgxpool.Pool, filepath string) error {
	script, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	if _, err := db.Exec(context.Background(), string(script)); err != nil {
		return err
	}
	return nil
}

func (pg *PGDB) SignUpUser(ctx context.Context, user *models.User) error {
	return pg.users.SignUpUser(ctx, user)
}

func (pg *PGDB) GetUser(ctx context.Context, login string) (*models.User, error) {
	return pg.users.GetUser(ctx, login)
}

func (pg *PGDB) GetAllUserItems(ctx context.Context, login string) ([]models.EncryptedItem, error) {
	return pg.items.GetAllUserItems(ctx, login)
}

func (pg *PGDB) GetUserItemsWithType(ctx context.Context, typ models.ItemType, login string) ([]models.EncryptedItem, error) {
	return pg.items.GetUserItemsWithType(ctx, typ, login)
}

func (pg *PGDB) AddItem(ctx context.Context, item *models.EncryptedItem) error {
	return pg.items.AddItem(ctx, item)
}

func (pg *PGDB) EditItem(ctx context.Context, item *models.EncryptedItem) error {
	return pg.items.EditItem(ctx, item)
}

func (pg *PGDB) DeleteItem(ctx context.Context, login string, itemID [16]byte) error {
	return pg.items.DeleteItem(ctx, login, itemID)
}

func (pg *PGDB) GetTypesCounts(ctx context.Context, login string) (map[models.ItemType]int32, error) {
	return pg.items.GetTypesCounts(ctx, login)
}
