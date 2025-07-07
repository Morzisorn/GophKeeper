package database

import (
	"context"
	"fmt"
	"gophkeeper/models"

	gen "gophkeeper/internal/server/repositories/database/generated"
)

type UserDatabase interface {
	SignUpUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, login string) (*models.User, error)
}

type UserDB struct {
	q *gen.Queries
}

func NewUserDB(q *gen.Queries) UserDatabase {
	return &UserDB{q: q}
}

func (db *UserDB) SignUpUser(ctx context.Context, user *models.User) error {
	return db.q.SignUpUser(ctx, gen.SignUpUserParams{
		Login:user.Login,
		Password: user.Password,
	})
}

func (db *UserDB) GetUser(ctx context.Context, login string) (*models.User, error) {
	user, err :=  db.q.GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("get user from db error: %w", err)
	}
	return &models.User{
		Login: user.Login,
		Password: user.Password,
	}, nil
}
