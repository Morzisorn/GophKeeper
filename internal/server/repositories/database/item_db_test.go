package database

import (
	"context"
	"fmt"
	gen "gophkeeper/internal/server/repositories/database/generated"
	"gophkeeper/models"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemDB_AddItem(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer func() {
		mock.Close()
	}()

	q := gen.New(mock)

	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	tests := []struct {
		name    string
		item    *models.EncryptedItem
		mockFn  func()
		wantErr bool
	}{
		{
			name: "successful add item",
			item: &models.EncryptedItem{
				UserLogin: "testuser",
				Name:      "test item",
				Type:      models.ItemType("CREDENTIALS"),
				EncryptedData: models.EncryptedData{
					EncryptedContent: "encrypted_content",
					Nonce:            "test_nonce",
				},
				Meta: models.Meta{},
			},
			mockFn: func() {
				mock.ExpectQuery("INSERT INTO items").
					WithArgs("testuser", "test item", itemTypeModelsToPg("CREDENTIALS"), "encrypted_content", "test_nonce", []byte(`{"Map":null}`)).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("550e8400-e29b-41d4-a716-446655440000"))
			},
			wantErr: false,
		},
		{
			name: "failed add item - database error",
			item: &models.EncryptedItem{
				UserLogin: "testuser",
				Name:      "test item",
				Type:      models.ItemTypeCREDENTIALS,
				EncryptedData: models.EncryptedData{
					EncryptedContent: "encrypted_content",
					Nonce:            "test_nonce",
				},
				Meta: models.Meta{Map: make(map[string]string)},
			},
			mockFn: func() {
				mock.ExpectQuery("INSERT INTO items").
					WithArgs("testuser", "test item", itemTypeModelsToPg("CREDENTIALS"), "encrypted_content", "test_nonce", []byte(`{"Map":null}`)).
					WillReturnError(fmt.Errorf("foreign key constraint fails"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := itemDB.AddItem(context.Background(), tt.item)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestItemDB_GetAllUserItems(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer func() {
		mock.Close()
	}()

	q := gen.New(mock)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	tests := []struct {
		name     string
		login    string
		mockFn   func()
		expected int // changed to check number of elements
		wantErr  bool
	}{
		{
			name:  "successful get all items",
			login: "testuser",
			mockFn: func() {
				// Use pgtype.UUID instead of [16]byte
				testUUID := pgtype.UUID{
					Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00},
					Valid: true,
				}

				rows := pgxmock.NewRows([]string{
					"id", "name", "type", "encrypted_data_content",
					"encrypted_data_nonce", "meta", "created_at", "updated_at",
				}).AddRow(
					testUUID, // use pgtype.UUID
					"test item",
					"CREDENTIALS",
					"encrypted_content",
					"test_nonce",
					[]byte(`{"Map":null}`),
					pgtype.Timestamp{Time: time.Now(), Valid: true},
					pgtype.Timestamp{Time: time.Now(), Valid: true},
				)
				mock.ExpectQuery("SELECT.*FROM items").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expected: 1,
			wantErr:  false,
		},
		{
			name:  "no items found",
			login: "emptyuser",
			mockFn: func() {
				rows := pgxmock.NewRows([]string{
					"id", "name", "type", "encrypted_data_content",
					"encrypted_data_nonce", "meta", "created_at", "updated_at",
				})
				mock.ExpectQuery("SELECT.*FROM items").
					WithArgs("emptyuser").
					WillReturnRows(rows)
			},
			expected: 0,
			wantErr:  false,
		},
		{
			name:  "database error",
			login: "erroruser",
			mockFn: func() {
				mock.ExpectQuery("SELECT.*FROM items").
					WithArgs("erroruser").
					WillReturnError(fmt.Errorf("database connection failed"))
			},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			items, err := itemDB.GetAllUserItems(context.Background(), tt.login)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, items, tt.expected)
				// Additional checks for successful case
				if tt.expected > 0 && len(items) > 0 {
					assert.Equal(t, "testuser", items[0].UserLogin)
					assert.Equal(t, "test item", items[0].Name)
					assert.Equal(t, models.ItemType("CREDENTIALS"), items[0].Type)
					assert.Equal(t, "encrypted_content", items[0].EncryptedData.EncryptedContent)
					assert.Equal(t, "test_nonce", items[0].EncryptedData.Nonce)
				}
			}
		})
	}
}

func TestItemDB_GetUserItemsWithType(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer func() {
		mock.Close()
	}()

	q := gen.New(mock)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	tests := []struct {
		name     string
		login    string
		itemType models.ItemType
		mockFn   func()
		expected int
		wantErr  bool
	}{
		{
			name:     "successful get items with type",
			login:    "testuser",
			itemType: models.ItemType("CREDENTIALS"),
			mockFn: func() {
				testUUID := pgtype.UUID{
					Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00},
					Valid: true,
				}
				rows := pgxmock.NewRows([]string{
					"id", "name", "type", "encrypted_data_content",
					"encrypted_data_nonce", "meta", "created_at", "updated_at",
				}).AddRow(
					testUUID,
					"login item",
					itemTypeModelsToPg(models.ItemTypeCREDENTIALS),
					"encrypted_content",
					"test_nonce",
					[]byte(`{"Map":null}`),
					pgtype.Timestamp{Time: time.Now(), Valid: true},
					pgtype.Timestamp{Time: time.Now(), Valid: true},
				)
				mock.ExpectQuery("SELECT.*FROM items.*WHERE.*type").
					WithArgs("testuser", itemTypeModelsToPg(models.ItemTypeCREDENTIALS)).
					WillReturnRows(rows)
			},
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "no items with type found",
			login:    "testuser",
			itemType: models.ItemTypeBINARY,
			mockFn: func() {
				rows := pgxmock.NewRows([]string{
					"id", "name", "type", "encrypted_data_content",
					"encrypted_data_nonce", "meta", "created_at", "updated_at",
				})
				mock.ExpectQuery("SELECT.*FROM items.*WHERE.*type").
					WithArgs("testuser", itemTypeModelsToPg(models.ItemTypeBINARY)).
					WillReturnRows(rows)
			},
			expected: 0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			items, err := itemDB.GetUserItemsWithType(context.Background(), tt.itemType, tt.login)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, items, tt.expected)
			}
		})
	}
}

func TestItemDB_GetTypesCounts(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer func() {
		mock.Close()
	}()

	q := gen.New(mock)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	tests := []struct {
		name     string
		login    string
		mockFn   func()
		expected map[models.ItemType]int32
		wantErr  bool
	}{
		{
			name:  "successful get types counts",
			login: "testuser",
			mockFn: func() {
				rows := pgxmock.NewRows([]string{"type", "count"}).
					AddRow(itemTypeModelsToPg(models.ItemTypeCREDENTIALS), int64(5)).
					AddRow(itemTypeModelsToPg(models.ItemTypeTEXT), int64(3)).
					AddRow(itemTypeModelsToPg(models.ItemTypeBINARY), int64(2))
				mock.ExpectQuery("SELECT.*type.*COUNT.*FROM items.*GROUP BY type").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expected: map[models.ItemType]int32{
				models.ItemTypeCREDENTIALS: 5,
				models.ItemTypeTEXT:        3,
				models.ItemTypeBINARY:      2,
			},
			wantErr: false,
		},
		{
			name:  "no items for user",
			login: "emptyuser",
			mockFn: func() {
				rows := pgxmock.NewRows([]string{"type", "count"})
				mock.ExpectQuery("SELECT.*type.*COUNT.*FROM items.*GROUP BY type").
					WithArgs("emptyuser").
					WillReturnRows(rows)
			},
			expected: map[models.ItemType]int32{},
			wantErr:  false,
		},
		{
			name:  "database error",
			login: "erroruser",
			mockFn: func() {
				mock.ExpectQuery("SELECT.*type.*COUNT.*FROM items.*GROUP BY type").
					WithArgs("erroruser").
					WillReturnError(fmt.Errorf("database error"))
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			counts, err := itemDB.GetTypesCounts(context.Background(), tt.login)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, counts)
			}
		})
	}
}

func TestItemDB_EditItem(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer func() {
		mock.Close()
	}()

	q := gen.New(mock)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	tests := []struct {
		name    string
		item    *models.EncryptedItem
		mockFn  func()
		wantErr bool
	}{
		{
			name: "successful edit item",
			item: &models.EncryptedItem{
				ID:        [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00},
				UserLogin: "testuser",
				Name:      "updated item",
				Type:      models.ItemTypeCREDENTIALS,
				EncryptedData: models.EncryptedData{
					EncryptedContent: "new_encrypted_content",
					Nonce:            "new_nonce",
				},
				Meta: models.Meta{},
			},
			mockFn: func() {
				mock.ExpectExec("UPDATE items SET").
					WithArgs(
						pgtype.UUID{Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00}, Valid: true},
						"updated item",
						"new_encrypted_content",
						"new_nonce",
						[]byte(`{"Map":null}`),
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))
			},
			wantErr: false,
		},
		{
			name: "failed edit item - item not found",
			item: &models.EncryptedItem{
				ID:        [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x01},
				UserLogin: "testuser",
				Name:      "nonexistent item",
				Type:      models.ItemTypeCREDENTIALS,
				EncryptedData: models.EncryptedData{
					EncryptedContent: "encrypted_content",
					Nonce:            "test_nonce",
				},
				Meta: models.Meta{},
			},
			mockFn: func() {
				mock.ExpectExec("UPDATE items SET").
					WithArgs(
						pgtype.UUID{Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x01}, Valid: true},
						"nonexistent item",
						"encrypted_content",
						"test_nonce",
						[]byte(`{"Map":null}`),
					).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))
			},
			wantErr: false, // UPDATE with 0 affected rows is not considered an error in this implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := itemDB.EditItem(context.Background(), tt.item)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestItemDB_DeleteItem(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer func() {
		mock.Close()
	}()

	q := gen.New(mock)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	tests := []struct {
		name    string
		login   string
		itemID  [16]byte
		mockFn  func()
		wantErr bool
	}{
		{
			name:   "successful delete item",
			login:  "testuser",
			itemID: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00},
			mockFn: func() {
				mock.ExpectExec("DELETE FROM items").
					WithArgs("testuser", pgtype.UUID{Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00}, Valid: true}).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
			},
			wantErr: false,
		},
		{
			name:   "item not found or unauthorized",
			login:  "testuser",
			itemID: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x01},
			mockFn: func() {
				mock.ExpectExec("DELETE FROM items").
					WithArgs("testuser", pgtype.UUID{Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x01}, Valid: true}).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))
			},
			wantErr: false, // DELETE with 0 affected rows is not considered an error in this implementation
		},
		{
			name:   "database error",
			login:  "testuser",
			itemID: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x02},
			mockFn: func() {
				mock.ExpectExec("DELETE FROM items").
					WithArgs("testuser", pgtype.UUID{Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x02}, Valid: true}).
					WillReturnError(fmt.Errorf("database connection failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := itemDB.DeleteItem(context.Background(), tt.login, tt.itemID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
