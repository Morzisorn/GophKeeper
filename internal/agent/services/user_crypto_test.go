package services

import (
	"github.com/stretchr/testify/require"
	"gophkeeper/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCryptoService(t *testing.T) {
	mockClient := &MockClient{}

	service, err := NewCryptoService(mockClient)
	require.NoError(t, err)

	assert.NotNil(t, service)
	assert.Equal(t, mockClient, service.Client)
	assert.NotNil(t, service.config)
}

func TestCryptoService_SetPublicKey_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.SetPublicKey()
	})
}

func TestCryptoService_SetPublicKey_NilClient(t *testing.T) {
	service := &CryptoService{
		Client: nil,
		config: &config.Config{},
	}

	assert.Panics(t, func() {
		service.SetPublicKey()
	})
}

func TestCryptoService_SetSalt_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.SetSalt("test-salt")
	})
}

func TestCryptoService_SetSalt_NilConfig(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: nil,
	}

	assert.Panics(t, func() {
		service.SetSalt("test-salt")
	})
}

func TestCryptoService_SetSalt_InvalidBase64(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: &config.Config{},
	}

	err := service.SetSalt("invalid-base64!")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "set salt error")
}

func TestCryptoService_GetSalt_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.GetSalt()
	})
}

func TestCryptoService_GetSalt_NilConfig(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: nil,
	}

	assert.Panics(t, func() {
		service.GetSalt()
	})
}

func TestCryptoService_SetMasterPassword_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.SetMasterPassword("test-password")
	})
}

func TestCryptoService_SetMasterPassword_NilConfig(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: nil,
	}

	assert.Panics(t, func() {
		service.SetMasterPassword("test-password")
	})
}

func TestCryptoService_GetMasterPassword_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.GetMasterPassword()
	})
}

func TestCryptoService_GetMasterPassword_NilConfig(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: nil,
	}

	assert.Panics(t, func() {
		service.GetMasterPassword()
	})
}

func TestCryptoService_GenerateMasterKey_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.GenerateMasterKey()
	})
}

func TestCryptoService_GenerateMasterKey_NilConfig(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: nil,
	}

	assert.Panics(t, func() {
		service.GenerateMasterKey()
	})
}

func TestCryptoService_SetMasterKey_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.SetMasterKey([]byte("test-key"))
	})
}

func TestCryptoService_SetMasterKey_NilConfig(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: nil,
	}

	assert.Panics(t, func() {
		service.SetMasterKey([]byte("test-key"))
	})
}

func TestCryptoService_GetMasterKey_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.GetMasterKey()
	})
}

func TestCryptoService_GetMasterKey_NilConfig(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: nil,
	}

	assert.Panics(t, func() {
		service.GetMasterKey()
	})
}

func TestCryptoService_SetGetMasterPassword_ValidService(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: &config.Config{
			AgentConfig: config.AgentConfig{},
		},
	}

	password := "test-master-password"
	service.SetMasterPassword(password)

	result := service.GetMasterPassword()
	assert.Equal(t, password, result)
}

func TestCryptoService_SetGetMasterKey_ValidService(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		config: &config.Config{
			AgentConfig: config.AgentConfig{},
		},
	}

	key := []byte("test-master-key-32-bytes-long!!")
	service.SetMasterKey(key)

	result := service.GetMasterKey()
	assert.Equal(t, key, result)
}

func TestEncryptData_NilPublicKey(t *testing.T) {
	// Этот тест потребует настройки config.GetAgentConfig()
	// Пропускаем, так как требует глобального состояния
	t.Skip("encryptData uses global config - needs integration test")
}
