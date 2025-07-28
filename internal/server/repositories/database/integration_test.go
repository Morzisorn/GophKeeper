package database

import (
	"context"
	"testing"
	"time"

	gen "gophkeeper/internal/server/repositories/database/generated"
	"gophkeeper/models"

	"github.com/jackc/pgx/v5/pgtype"
	pgxmock "github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationUserItemOperations(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := gen.New(mock)

	userDB, err := NewUserDB(q, mock)
	require.NoError(t, err)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	pgdb := &PGDB{
		users: userDB,
		items: itemDB,
	}

	ctx := context.Background()

	// Test user signup
	user := &models.User{
		Login:    "integrationuser",
		Password: []byte("hashedpassword"),
		Salt:     "testsalt",
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs("integrationuser", []byte("hashedpassword"), "testsalt").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = pgdb.SignUpUser(ctx, user)
	assert.NoError(t, err)

	// Test get user
	mock.ExpectQuery("SELECT login, password, salt FROM users").
		WithArgs("integrationuser").
		WillReturnRows(pgxmock.NewRows([]string{"login", "password", "salt"}).
			AddRow("integrationuser", []byte("hashedpassword"), "testsalt"))

	retrievedUser, err := pgdb.GetUser(ctx, "integrationuser")
	assert.NoError(t, err)
	assert.Equal(t, user.Login, retrievedUser.Login)
	assert.Equal(t, user.Password, retrievedUser.Password)
	assert.Equal(t, user.Salt, retrievedUser.Salt)

	// Test add item
	item := &models.EncryptedItem{
		UserLogin: "integrationuser",
		Name:      "test credential",
		Type:      models.ItemTypeCREDENTIALS,
		EncryptedData: models.EncryptedData{
			EncryptedContent: "encrypted_login_password",
			Nonce:            "random_nonce",
		},
		Meta: models.Meta{},
	}

	// Используем правильный JSON для Meta
	mock.ExpectQuery("INSERT INTO items").
		WithArgs("integrationuser", "test credential", itemTypeModelsToPg(models.ItemTypeCREDENTIALS), "encrypted_login_password", "random_nonce", []byte(`{"Map":null}`)).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow("550e8400-e29b-41d4-a716-446655440000"))

	err = pgdb.AddItem(ctx, item)
	assert.NoError(t, err)

	// Test get all user items - исправляем порядок ожиданий
	testUUID := pgtype.UUID{
		Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00},
		Valid: true,
	}

	mock.ExpectQuery("SELECT.*FROM items").
		WithArgs("integrationuser").
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "name", "type", "encrypted_data_content",
			"encrypted_data_nonce", "meta", "created_at", "updated_at",
		}).AddRow(
			testUUID,
			"test credential",
			itemTypeModelsToPg(models.ItemTypeCREDENTIALS),
			"encrypted_login_password",
			"random_nonce",
			[]byte(`{"Map":null}`),
			pgtype.Timestamp{Time: time.Now(), Valid: true},
			pgtype.Timestamp{Time: time.Now(), Valid: true},
		))

	items, err := pgdb.GetAllUserItems(ctx, "integrationuser")
	assert.NoError(t, err)
	require.Len(t, items, 1) // используем require чтобы избежать panic
	if len(items) > 0 {
		assert.Equal(t, "test credential", items[0].Name)
		assert.Equal(t, models.ItemTypeCREDENTIALS, items[0].Type)
	}

	// Test get types counts
	mock.ExpectQuery("SELECT.*type.*COUNT.*FROM items.*GROUP BY type").
		WithArgs("integrationuser").
		WillReturnRows(pgxmock.NewRows([]string{"type", "count"}).
			AddRow(itemTypeModelsToPg(models.ItemTypeCREDENTIALS), int64(1)))

	counts, err := pgdb.GetTypesCounts(ctx, "integrationuser")
	assert.NoError(t, err)
	assert.Equal(t, int32(1), counts[models.ItemTypeCREDENTIALS])
}
