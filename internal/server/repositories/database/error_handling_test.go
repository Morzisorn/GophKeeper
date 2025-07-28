package database

import (
	"context"
	"fmt"
	"testing"
	"time"

	gen "gophkeeper/internal/server/repositories/database/generated"
	"gophkeeper/models"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandling_InvalidJSONMeta(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := gen.New(mock)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	// Test GetAllUserItems with invalid JSON meta
	testUUID := pgtype.UUID{
		Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00},
		Valid: true,
	}

	mock.ExpectQuery("SELECT.*FROM items").
		WithArgs("testuser").
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "name", "type", "encrypted_data_content",
			"encrypted_data_nonce", "meta", "created_at", "updated_at",
		}).AddRow(
			testUUID,
			"test item",
			"CREDENTIALS",
			"encrypted_content",
			"test_nonce",
			[]byte("invalid json"),
			pgtype.Timestamp{Time: time.Now(), Valid: true},
			pgtype.Timestamp{Time: time.Now(), Valid: true},
		))

	items, err := itemDB.GetAllUserItems(context.Background(), "testuser")
	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Contains(t, err.Error(), "unmarshal meta info error")
}

func TestErrorHandling_DatabaseConnectionFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := gen.New(mock)
	userDB, err := NewUserDB(q, mock)
	require.NoError(t, err)

	// Test connection failure during user signup
	mock.ExpectExec("INSERT INTO users").
		WithArgs("testuser", []byte("password"), "salt").
		WillReturnError(fmt.Errorf("connection lost"))

	user := &models.User{
		Login:    "testuser",
		Password: []byte("password"),
		Salt:     "salt",
	}

	err = userDB.SignUpUser(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection lost")
}

func TestErrorHandling_QueryScanFailure(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := gen.New(mock)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	// Test scan failure during GetAllUserItems
	rows := pgxmock.NewRows([]string{
		"id", "name", "type", "encrypted_data_content",
		"encrypted_data_nonce", "meta", "created_at", "updated_at",
	}).AddRow(
		"invalid_uuid_format", // This will cause scan failure
		"test item",
		"CREDENTIALS",
		"encrypted_content",
		"test_nonce",
		[]byte(`{"Map":null}`),
		pgtype.Timestamp{Time: time.Now(), Valid: true},
		pgtype.Timestamp{Time: time.Now(), Valid: true},
	).RowError(0, fmt.Errorf("scan error"))

	mock.ExpectQuery("SELECT.*FROM items").
		WithArgs("testuser").
		WillReturnRows(rows)

	items, err := itemDB.GetAllUserItems(context.Background(), "testuser")
	assert.Error(t, err)
	assert.Nil(t, items)
}
