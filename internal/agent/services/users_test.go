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
	cnfg, err := config.NewAgentConfig()
	require.NoError(t, err)

	mockClient := &MockClient{}
	mockCrypto := &CryptoService{}

	service, err := NewUserService(cnfg, mockClient, mockCrypto)
	require.NoError(t, err)

	assert.NotNil(t, service)
	assert.Equal(t, mockClient, service.Client)
	assert.Equal(t, mockCrypto, service.crypto)
	assert.NotNil(t, service.cnfg)
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
	cnfg, err := config.NewAgentConfig()
	require.NoError(t, err)

	service := &UserService{
		Client: &MockClient{},
		crypto: &CryptoService{},
		cnfg:   cnfg,
	}

	user := &models.User{
		Login:    "",
		Password: []byte("test-password"),
	}

	err = service.SignUpUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrRequiredArgumentIsMissing, err)
}

func TestUserService_SignUpUser_NilPassword(t *testing.T) {
	cnfg, err := config.NewAgentConfig()
	require.NoError(t, err)

	service := &UserService{
		Client: &MockClient{},
		crypto: &CryptoService{},
		cnfg:   cnfg,
	}

	user := &models.User{
		Login:    "test-login",
		Password: nil,
	}

	err = service.SignUpUser(context.Background(), user)

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
	cnfg, err := config.NewAgentConfig()
	require.NoError(t, err)

	service := &UserService{
		Client: &MockClient{},
		crypto: &CryptoService{},
		cnfg:   cnfg,
	}

	user := &models.User{
		Login:    "",
		Password: []byte("test-password"),
	}

	err = service.SignInUser(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, errs.ErrRequiredArgumentIsMissing, err)
}

func TestUserService_SignInUser_NilPassword(t *testing.T) {
	cnfg, err := config.NewAgentConfig()
	require.NoError(t, err)

	service := &UserService{
		Client: &MockClient{},
		crypto: &CryptoService{},
		cnfg:   cnfg,
	}

	user := &models.User{
		Login:    "test-login",
		Password: nil,
	}

	err = service.SignInUser(context.Background(), user)

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
	cnfg, err := config.NewAgentConfig()
	require.NoError(t, err)

	service := &UserService{
		Client: &MockClient{},
		crypto: nil,
		cnfg:   cnfg,
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
		crypto: &CryptoService{},
		cnfg:   nil,
	}

	assert.Panics(t, func() {
		service.Logout()
	})
}
