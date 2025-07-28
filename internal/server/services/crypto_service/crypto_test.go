package crypto_service

import (
	"context"
	"gophkeeper/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStorage struct{}

func (m *MockStorage) SignUpUser(ctx context.Context, user *models.User) error { return nil }
func (m *MockStorage) GetUser(ctx context.Context, login string) (*models.User, error) {
	return nil, nil
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

func TestNewCryptoService(t *testing.T) {
	repo := &MockStorage{}
	service, err := NewCryptoService(repo)

	assert.NoError(t, err)
	assert.NotNil(t, service)
}
