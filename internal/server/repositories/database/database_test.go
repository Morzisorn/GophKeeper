package database

import (
	"context"
	"fmt"
	"testing"

	"gophkeeper/models"

	gen "gophkeeper/internal/server/repositories/database/generated"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPGDB_SignUpUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer func() {
		mock.Close()
	}()

	q := gen.New(mock)

	userDB, err := NewUserDB(q, mock)
	require.NoError(t, err)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	pgdb := &PGDB{
		users: userDB,
		items: itemDB,
	}

	tests := []struct {
		name    string
		user    *models.User
		mockFn  func()
		wantErr bool
	}{
		{
			name: "successful user signup",
			user: &models.User{
				Login:    "testuser",
				Password: []byte("hashedpassword"),
				Salt:     "testsalt",
			},
			mockFn: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("testuser", []byte("hashedpassword"), "testsalt").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: false,
		},
		{
			name: "failed user signup - duplicate login",
			user: &models.User{
				Login:    "existinguser",
				Password: []byte("hashedpassword"),
				Salt:     "testsalt",
			},
			mockFn: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("existinguser", []byte("hashedpassword"), "testsalt").
					WillReturnError(fmt.Errorf("duplicate key value violates unique constraint"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := pgdb.SignUpUser(context.Background(), tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPGDB_GetUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer func() {
		mock.Close()
	}()

	q := gen.New(mock)

	userDB, err := NewUserDB(q, mock)
	require.NoError(t, err)
	itemDB, err := NewItemDB(q, mock)
	require.NoError(t, err)

	pgdb := &PGDB{
		users: userDB,
		items: itemDB,
	}

	tests := []struct {
		name     string
		login    string
		mockFn   func()
		expected *models.User
		wantErr  bool
	}{
		{
			name:  "successful get user",
			login: "testuser",
			mockFn: func() {
				rows := pgxmock.NewRows([]string{"login", "password", "salt"}).
					AddRow("testuser", []byte("hashedpassword"), "testsalt")
				mock.ExpectQuery("SELECT login, password, salt FROM users").
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			expected: &models.User{
				Login:    "testuser",
				Password: []byte("hashedpassword"),
				Salt:     "testsalt",
			},
			wantErr: false,
		},
		{
			name:  "user not found",
			login: "nonexistent",
			mockFn: func() {
				mock.ExpectQuery("SELECT login, password, salt FROM users").
					WithArgs("nonexistent").
					WillReturnError(pgx.ErrNoRows)
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			user, err := pgdb.GetUser(context.Background(), tt.login)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}
		})
	}
}

// Test helper functions
func TestItemTypeModelsToPg(t *testing.T) {
	tests := []struct {
		name     string
		itemType models.ItemType
		expected string
	}{
		{
			name:     "credentials type",
			itemType: models.ItemType("CREDENTIALS"),
			expected: "CREDENTIALS",
		},
		{
			name:     "text type",
			itemType: models.ItemType("TEXT"),
			expected: "TEXT",
		},
		{
			name:     "binary type",
			itemType: models.ItemType("BINARY"),
			expected: "BINARY",
		},
		{
			name:     "card type",
			itemType: models.ItemType("CARD"),
			expected: "CARD",
		},
		{
			name:     "unknown type",
			itemType: models.ItemType("UNKNOWN"),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := itemTypeModelsToPg(tt.itemType)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}
