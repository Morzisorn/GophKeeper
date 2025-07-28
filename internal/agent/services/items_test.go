package services

import (
	"context"
	"github.com/stretchr/testify/require"
	"gophkeeper/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewItemService(t *testing.T) {
	mockClient := &MockClient{}
	mockCrypto := &CryptoService{}

	service, err := NewItemService(mockClient, mockCrypto)
	require.NoError(t, err)

	assert.NotNil(t, service)
	assert.Equal(t, mockClient, service.Client)
	assert.Equal(t, mockCrypto, service.Crypto)
}

func TestItemService_AddItem_NilService(t *testing.T) {
	var service *ItemService = nil

	item := &models.Item{}

	assert.Panics(t, func() {
		service.AddItem(context.Background(), item)
	})
}

func TestItemService_AddItem_NilCrypto(t *testing.T) {
	service := &ItemService{
		Client: &MockClient{},
		Crypto: nil,
	}

	item := &models.Item{}

	assert.Panics(t, func() {
		service.AddItem(context.Background(), item)
	})
}

func TestItemService_EditItem_NilService(t *testing.T) {
	var service *ItemService = nil

	item := &models.Item{}

	assert.Panics(t, func() {
		service.EditItem(context.Background(), item)
	})
}

func TestItemService_EditItem_NilCrypto(t *testing.T) {
	service := &ItemService{
		Client: &MockClient{},
		Crypto: nil,
	}

	item := &models.Item{}

	assert.Panics(t, func() {
		service.EditItem(context.Background(), item)
	})
}

func TestItemService_DeleteItem_NilService(t *testing.T) {
	var service *ItemService = nil

	var itemID [16]byte

	assert.Panics(t, func() {
		service.DeleteItem(context.Background(), "test-login", itemID)
	})
}

func TestItemService_DeleteItem_NilClient(t *testing.T) {
	service := &ItemService{
		Client: nil,
		Crypto: &CryptoService{},
	}

	var itemID [16]byte

	assert.Panics(t, func() {
		service.DeleteItem(context.Background(), "test-login", itemID)
	})
}

func TestItemService_GetItems_NilService(t *testing.T) {
	var service *ItemService = nil

	assert.Panics(t, func() {
		service.GetItems(context.Background(), "test-login", models.ItemTypeCARD)
	})
}

func TestItemService_GetItems_NilClient(t *testing.T) {
	service := &ItemService{
		Client: nil,
		Crypto: &CryptoService{},
	}

	assert.Panics(t, func() {
		service.GetItems(context.Background(), "test-login", models.ItemTypeCARD)
	})
}

func TestItemService_GetTypesCounts_NilService(t *testing.T) {
	var service *ItemService = nil

	assert.Panics(t, func() {
		service.GetTypesCounts(context.Background(), "test-login")
	})
}

func TestItemService_GetTypesCounts_NilClient(t *testing.T) {
	service := &ItemService{
		Client: nil,
		Crypto: &CryptoService{},
	}

	assert.Panics(t, func() {
		service.GetTypesCounts(context.Background(), "test-login")
	})
}

func TestItemService_DecryptItem_NilService(t *testing.T) {
	var service *ItemService = nil

	encItem := &models.EncryptedItem{}

	assert.Panics(t, func() {
		service.DecryptItem(encItem)
	})
}

func TestItemService_DecryptItem_NilCrypto(t *testing.T) {
	service := &ItemService{
		Client: &MockClient{},
		Crypto: nil,
	}

	encItem := &models.EncryptedItem{}

	assert.Panics(t, func() {
		service.DecryptItem(encItem)
	})
}

// MockClient для тестирования
type MockClient struct{}

func (m *MockClient) SignUpUser(ctx context.Context, user *models.User) (token string, salt string, err error) {
	return "", "", nil
}

func (m *MockClient) SignInUser(ctx context.Context, user *models.User) (token string, salt string, err error) {
	return "", "", nil
}

func (m *MockClient) SetJWTToken(token string) {}

func (m *MockClient) GetJWTToken() string {
	return ""
}

func (m *MockClient) GetPublicKeyPEM(ctx context.Context) (string, error) {
	return "", nil
}

func (m *MockClient) AddItem(ctx context.Context, item *models.EncryptedItem) error {
	return nil
}

func (m *MockClient) EditItem(ctx context.Context, item *models.EncryptedItem) error {
	return nil
}

func (m *MockClient) DeleteItem(ctx context.Context, login string, itemID [16]byte) error {
	return nil
}

func (m *MockClient) GetItems(ctx context.Context, login string, typ models.ItemType) ([]models.EncryptedItem, error) {
	return nil, nil
}

func (m *MockClient) GetTypesCounts(ctx context.Context, login string) (map[string]int32, error) {
	return nil, nil
}
