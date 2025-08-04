package services

import (
	"gophkeeper/config"
	"gophkeeper/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCryptoService_EncryptItem_EmptyMasterPassword(t *testing.T) {
	cnfg, err := config.NewAgentConfig()
	assert.NoError(t, err)

	cryptoService, err := NewCryptoService(cnfg, nil)
	assert.NoError(t, err)

	itemService, err := NewItemService(nil, cryptoService)
	assert.NoError(t, err)

	item := &models.Item{}

	err = itemService.AddItem(nil, item)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "master password")
}

func TestNewCryptoService(t *testing.T) {
	cnfg, err := config.NewAgentConfig()
	assert.NoError(t, err)

	service, err := NewCryptoService(cnfg, nil)
	assert.NoError(t, err)
	assert.NotNil(t, service)
}

func TestNewItemService(t *testing.T) {
	cnfg, err := config.NewAgentConfig()
	assert.NoError(t, err)

	cryptoService, err := NewCryptoService(cnfg, nil)
	assert.NoError(t, err)

	itemService, err := NewItemService(nil, cryptoService)
	assert.NoError(t, err)
	assert.NotNil(t, itemService)
}
