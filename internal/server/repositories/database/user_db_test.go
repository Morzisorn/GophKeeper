package database

import (
	"context"
	"fmt"
	gen "gophkeeper/internal/server/repositories/database/generated"
	"gophkeeper/models"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// user_db_test.go
func TestUserDB_SignUpUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := gen.New(mock)

	userDB, err := NewUserDB(q, mock)
	require.NoError(t, err)

	tests := []struct {
		name    string
		user    *models.User
		mockFn  func()
		wantErr bool
	}{
		{
			name: "successful signup",
			user: &models.User{
				Login:    "newuser",
				Password: []byte("password123"),
				Salt:     "randomsalt",
			},
			mockFn: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("newuser", []byte("password123"), "randomsalt").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
			wantErr: false,
		},
		{
			name: "failed signup - database error",
			user: &models.User{
				Login:    "testuser",
				Password: []byte("password123"),
				Salt:     "randomsalt",
			},
			mockFn: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("testuser", []byte("password123"), "randomsalt").
					WillReturnError(fmt.Errorf("connection error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := userDB.SignUpUser(context.Background(), tt.user)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserDB_GetUser(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	q := gen.New(mock)

	userDB, err := NewUserDB(q, mock)
	require.NoError(t, err)

	tests := []struct {
		name     string
		login    string
		mockFn   func()
		expected *models.User
		wantErr  bool
	}{
		{
			name:  "successful get user",
			login: "existinguser",
			mockFn: func() {
				rows := pgxmock.NewRows([]string{"login", "password", "salt"}).
					AddRow("existinguser", []byte("hashedpass"), "usersalt")
				mock.ExpectQuery("SELECT login, password, salt FROM users").
					WithArgs("existinguser").
					WillReturnRows(rows)
			},
			expected: &models.User{
				Login:    "existinguser",
				Password: []byte("hashedpass"),
				Salt:     "usersalt",
			},
			wantErr: false,
		},
		{
			name:  "user not found",
			login: "notfound",
			mockFn: func() {
				mock.ExpectQuery("SELECT login, password, salt FROM users").
					WithArgs("notfound").
					WillReturnError(pgx.ErrNoRows)
			},
			expected: nil,
			wantErr:  true,
		},
		{
			name:  "database error",
			login: "erroruser",
			mockFn: func() {
				mock.ExpectQuery("SELECT login, password, salt FROM users").
					WithArgs("erroruser").
					WillReturnError(fmt.Errorf("database connection failed"))
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			user, err := userDB.GetUser(context.Background(), tt.login)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, user)
			}
		})
	}
}
