package services

import (
	"context"
	"github.com/stretchr/testify/require"
	"gophkeeper/config"
	"gophkeeper/internal/errs"
	"gophkeeper/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserService(t *testing.T) {
	mockClient := &MockClient{}
	mockCrypto := &CryptoService{}

	service, err := NewUserService(mockClient, mockCrypto)
	require.NoError(t, err)

	assert.NotNil(t, service)
	assert.Equal(t, mockClient, service.Client)
	assert.Equal(t, mockCrypto, service.Crypto)
	assert.NotNil(t, service.config)
}

func TestUserService_SignUpUser_NilService(t *testing.T) {
	var service *UserService = nil

	user := &models.User{
		Login:    "test-login",
		Password: []byte("test-password"),
	}

	// Этот вызов должен привести к panic при попытке вызвать encryptData
	// или обратиться к service.Client
	assert.Panics(t, func() {
		service.SignUpUser(context.Background(), user)
	})
}

func TestUserService_SignUpUser_EmptyLogin(t *testing.T) {
	service := &UserService{
		Client: &MockClient{},
		Crypto: &CryptoService{},
		config: &config.Config{},
	}

	user := &models.User{
		Login:    "",
		Password: []byte("test-password"),
	}

	err := service.SignUpUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrRequiredArgumentIsMissing, err)
}

func TestUserService_SignUpUser_NilPassword(t *testing.T) {
	service := &UserService{
		Client: &MockClient{},
		Crypto: &CryptoService{},
		config: &config.Config{},
	}

	user := &models.User{
		Login:    "test-login",
		Password: nil,
	}

	err := service.SignUpUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrRequiredArgumentIsMissing, err)
}

func TestUserService_SignInUser_NilService(t *testing.T) {
	var service *UserService = nil

	user := &models.User{
		Login:    "test-login",
		Password: []byte("test-password"),
	}

	// Этот вызов должен привести к panic при попытке вызвать encryptData
	// или обратиться к service.Client
	assert.Panics(t, func() {
		service.SignInUser(context.Background(), user)
	})
}

func TestUserService_SignInUser_EmptyLogin(t *testing.T) {
	service := &UserService{
		Client: &MockClient{},
		Crypto: &CryptoService{},
		config: &config.Config{},
	}

	user := &models.User{
		Login:    "",
		Password: []byte("test-password"),
	}

	err := service.SignInUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrRequiredArgumentIsMissing, err)
}

func TestUserService_SignInUser_NilPassword(t *testing.T) {
	service := &UserService{
		Client: &MockClient{},
		Crypto: &CryptoService{},
		config: &config.Config{},
	}

	user := &models.User{
		Login:    "test-login",
		Password: nil,
	}

	err := service.SignInUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrRequiredArgumentIsMissing, err)
}

func TestUserService_SetMasterKey_NilService(t *testing.T) {
	var service *UserService = nil

	assert.Panics(t, func() {
		service.SetMasterKey("test-password")
	})
}

func TestUserService_SetMasterKey_NilCrypto(t *testing.T) {
	service := &UserService{
		Client: &MockClient{},
		Crypto: nil,
		config: &config.Config{},
	}

	assert.Panics(t, func() {
		service.SetMasterKey("test-password")
	})
}

func TestUserService_Logout_NilService(t *testing.T) {
	var service *UserService = nil

	assert.Panics(t, func() {
		service.Logout()
	})
}

func TestUserService_Logout_NilConfig(t *testing.T) {
	service := &UserService{
		Client: &MockClient{},
		Crypto: &CryptoService{},
		config: nil,
	}

	assert.Panics(t, func() {
		service.Logout()
	})
}

func TestUserService_Logout_ValidService(t *testing.T) {
	service := &UserService{
		Client: &MockClient{},
		Crypto: &CryptoService{},
		config: &config.Config{
			AgentConfig: config.AgentConfig{
				MasterKey:      []byte("test-key"),
				MasterPassword: "test-password",
				Salt:           []byte("test-salt"),
			},
		},
	}

	err := service.Logout()

	assert.NoError(t, err)
	assert.Nil(t, service.config.MasterKey)
	assert.Empty(t, service.config.MasterPassword)
	assert.Nil(t, service.config.Salt)
}
