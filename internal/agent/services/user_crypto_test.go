package services

import (
	"github.com/stretchr/testify/require"
	"gophkeeper/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoService_SetPublicKey_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.SetPublicKey()
	})
}

func TestCryptoService_SetPublicKey_NilClient(t *testing.T) {
	cnfg, err := config.NewAgentConfig()
	require.NoError(t, err)

	service := &CryptoService{
		Client: nil,
		cnfg:   cnfg,
	}

	assert.Panics(t, func() {
		service.SetPublicKey()
	})
}

func TestCryptoService_SetSalt_InvalidBase64(t *testing.T) {
	cnfg, err := config.NewAgentConfig()
	require.NoError(t, err)

	service := &CryptoService{
		Client: &MockClient{},
		cnfg:   cnfg,
	}

	err = service.setSalt("invalid-base64!")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "set salt error")
}

func TestCryptoService_SetMasterPassword_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.setMasterPassword("test-password")
	})
}

func TestCryptoService_SetMasterPassword_NilConfig(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		cnfg:   nil,
	}

	assert.Panics(t, func() {
		service.setMasterPassword("test-password")
	})
}

func TestCryptoService_GenerateMasterKey_NilService(t *testing.T) {
	var service *CryptoService = nil

	assert.Panics(t, func() {
		service.generateMasterKey()
	})
}

func TestCryptoService_GenerateMasterKey_NilConfig(t *testing.T) {
	service := &CryptoService{
		Client: &MockClient{},
		cnfg:   nil,
	}

	assert.Panics(t, func() {
		service.generateMasterKey()
	})
}

func TestEncryptData_NilPublicKey(t *testing.T) {
	// This test would require config.GetAgentConfig() setup
	// Skip since it requires global state
	t.Skip("encryptData uses global config - needs integration test")
}
