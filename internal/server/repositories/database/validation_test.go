package database

import (
	"context"
	"fmt"
	"testing"

	gen "gophkeeper/internal/server/repositories/database/generated"
	"gophkeeper/models"

	pgxmock "github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidation_EmptyUserLogin(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := gen.New(mock)

	userDB := NewUserDB(q, mock)

	// Test with empty login
	mock.ExpectExec("INSERT INTO users").
		WithArgs("", []byte("password"), "salt").
		WillReturnError(fmt.Errorf("check constraint violation"))

	user := &models.User{
		Login:    "",
		Password: []byte("password"),
		Salt:     "salt",
	}

	err = userDB.SignUpUser(context.Background(), user)
	assert.Error(t, err)
}

func TestValidation_NilPassword(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := gen.New(mock)

	userDB := NewUserDB(q, mock)

	// Test with nil password
	mock.ExpectExec("INSERT INTO users").
		WithArgs("testuser", nil, "salt").
		WillReturnError(fmt.Errorf("null value in column violates not-null constraint"))

	user := &models.User{
		Login:    "testuser",
		Password: nil,
		Salt:     "salt",
	}

	err = userDB.SignUpUser(context.Background(), user)
	assert.Error(t, err)
}

func TestValidation_EmptyItemName(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := gen.New(mock)

	itemDB := NewItemDB(q, mock)

	// Test with empty item name
	mock.ExpectQuery("INSERT INTO items").
		WithArgs("testuser", "", "CREDENTIALS", "content", "nonce", []byte("{}")).
		WillReturnError(fmt.Errorf("check constraint violation"))

	item := &models.EncryptedItem{
		UserLogin: "testuser",
		Name:      "",
		Type:      models.ItemType("CREDENTIALS"),
		EncryptedData: models.EncryptedData{
			EncryptedContent: "content",
			Nonce:            "nonce",
		},
		Meta: models.Meta{},
	}

	err = itemDB.AddItem(context.Background(), item)
	assert.Error(t, err)
}
