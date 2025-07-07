package database

import (
	"context"
	"fmt"
	"gophkeeper/config"
	"gophkeeper/models"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	gen "gophkeeper/internal/server/repositories/database/generated"
)

type Database interface {
	UserDatabase
	CredentialsDatabase
}

type PGDB struct {
	users UserDatabase
	creds CredentialsDatabase
}

func NewPGDB(cfg *config.Config) (Database, error) {
	pool, err := pgxpool.New(context.Background(), cfg.DBConnStr)
	if err != nil {
		return nil, fmt.Errorf("create new db error: %v", err)
	}

	err = createTables(pool)
	if err != nil {
		return nil, fmt.Errorf("create db tables error: %v", err)
	}
	
	q := gen.New(pool) 

	return &PGDB{
		users: NewUserDB(q),
		creds: NewCredsDB(q),
	}, nil
}

func createTables(db *pgxpool.Pool) error {
	rootDir, err := config.GetProjectRoot()
	if err != nil {
		return err
	}
	filepath := filepath.Join(rootDir, "internal", "server", "repositories", "database", "schema", "schema.sql")

	script, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), string(script))
	if err != nil {
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

// func (db *PGDB) Close() error {
	
// }