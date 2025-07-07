package database

import (
	gen "gophkeeper/internal/server/repositories/database/generated"
)

type CredentialsDatabase interface {
}

type CredentialsDB struct {
	q *gen.Queries
}

func NewCredsDB(q *gen.Queries) CredentialsDatabase {
	return &CredentialsDB{q: q}
}
