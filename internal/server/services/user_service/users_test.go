package user_service

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/config"
	"gophkeeper/internal/errs"
	"gophkeeper/models"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

// MockStorage реализует интерфейс repositories.Storage для тестирования
type MockStorage struct {
	shouldFail bool
	users      map[string]*models.User
}

func (m *MockStorage) SignUpUser(ctx context.Context, user *models.User) error {
	if m.shouldFail {
		return errors.New("storage error")
	}
	if m.users == nil {
		m.users = make(map[string]*models.User)
	}
	m.users[user.Login] = user
	return nil
}

func (m *MockStorage) GetUser(ctx context.Context, login string) (*models.User, error) {
	if m.shouldFail {
		return nil, errors.New("storage error")
	}
	if m.users == nil {
		m.users = make(map[string]*models.User)
	}
	user, exists := m.users[login]
	if !exists {
		return nil, pgx.ErrNoRows
	}
	return user, nil
}

func (m *MockStorage) GetAllUserItems(ctx context.Context, login string) ([]models.EncryptedItem, error) {
	return nil, nil
}
func (m *MockStorage) GetUserItemsWithType(ctx context.Context, typ models.ItemType, login string) ([]models.EncryptedItem, error) {
	return nil, nil
}
func (m *MockStorage) AddItem(ctx context.Context, item *models.EncryptedItem) error  { return nil }
func (m *MockStorage) EditItem(ctx context.Context, item *models.EncryptedItem) error { return nil }
func (m *MockStorage) DeleteItem(ctx context.Context, login string, itemID [16]byte) error {
	return nil
}
func (m *MockStorage) GetTypesCounts(ctx context.Context, login string) (map[models.ItemType]int32, error) {
	return nil, nil
}

func TestNewUserService(t *testing.T) {
	cnfg, err := config.NewServerConfig()
	assert.NoError(t, err)

	repo := &MockStorage{}
	service, err := NewUserService(cnfg, repo)
	assert.NoError(t, err)

	assert.NotNil(t, service)
}

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		setup   func(*MockStorage)
		wantErr bool
	}{
		{
			name:  "successful get user",
			login: "existinguser",
			setup: func(m *MockStorage) {
				m.users = map[string]*models.User{
					"existinguser": {
						Login:    "existinguser",
						Password: []byte("password"),
						Salt:     "salt",
					},
				}
			},
			wantErr: false,
		},
		{
			name:  "user not found",
			login: "nonexistent",
			setup: func(m *MockStorage) {
				m.users = make(map[string]*models.User)
			},
			wantErr: true,
		},
		{
			name:  "storage error",
			login: "testuser",
			setup: func(m *MockStorage) {
				m.shouldFail = true
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStorage{}
			tt.setup(mockRepo)
			cnfg, err := config.NewServerConfig()
			assert.NoError(t, err)
			service, err := NewUserService(cnfg, mockRepo)
			assert.NoError(t, err)

			user, err := service.GetUser(context.Background(), &models.User{Login: tt.login})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.login, user.Login)
			}
		})
	}
}

func TestUserService_SignUpUser_ValidationOnly(t *testing.T) {
	// Тестируем только логику валидации без шифрования
	tests := []struct {
		name              string
		login             string
		encryptedPassword string
		setup             func(*MockStorage)
		wantErr           bool
		expectedErrType   error
	}{
		{
			name:              "user already exists",
			login:             "existinguser",
			encryptedPassword: "dummy_encrypted",
			setup: func(m *MockStorage) {
				m.users = map[string]*models.User{
					"existinguser": {
						Login:    "existinguser",
						Password: []byte("password"),
						Salt:     "salt",
					},
				}
			},
			wantErr:         true,
			expectedErrType: errs.ErrUserAlreadyRegistered,
		},
		{
			name:              "storage error during user check",
			login:             "newuser",
			encryptedPassword: "dummy_encrypted",
			setup: func(m *MockStorage) {
				m.shouldFail = true
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStorage{}
			tt.setup(mockRepo)
			cnfg, err := config.NewServerConfig()
			assert.NoError(t, err)
			service, err := NewUserService(cnfg, mockRepo)
			assert.NoError(t, err)

			token, salt, err := service.SignUpUser(context.Background(), tt.login, tt.encryptedPassword)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Empty(t, salt)

				if tt.expectedErrType != nil {
					assert.ErrorIs(t, err, tt.expectedErrType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, salt)
			}
		})
	}
}

func TestUserService_SignInUser_ValidationOnly(t *testing.T) {
	// Тестируем только логику валидации без шифрования
	tests := []struct {
		name              string
		login             string
		encryptedPassword string
		setup             func(*MockStorage)
		wantErr           bool
		expectedErrType   error
	}{
		{
			name:              "user not found",
			login:             "nonexistent",
			encryptedPassword: "dummy_encrypted",
			setup: func(m *MockStorage) {
				m.users = make(map[string]*models.User)
			},
			wantErr:         true,
			expectedErrType: errs.ErrUserNotFound,
		},
		{
			name:              "storage error during user get",
			login:             "testuser",
			encryptedPassword: "dummy_encrypted",
			setup: func(m *MockStorage) {
				m.shouldFail = true
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStorage{}
			tt.setup(mockRepo)
			cnfg, err := config.NewServerConfig()
			assert.NoError(t, err)
			service, err := NewUserService(cnfg, mockRepo)
			assert.NoError(t, err)

			token, salt, err := service.SignInUser(context.Background(), tt.login, tt.encryptedPassword)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Empty(t, salt)

				if tt.expectedErrType != nil {
					assert.ErrorIs(t, err, tt.expectedErrType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.NotEmpty(t, salt)
			}
		})
	}
}

func TestDecryptPassword_InvalidInput(t *testing.T) {
	t.Skip("Skipping decryptPassword tests - requires RSA private key configuration")
}
