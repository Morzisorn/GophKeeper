package item_service

import (
	"context"
	"errors"
	"testing"

	"gophkeeper/models"

	"github.com/stretchr/testify/assert"
)

// MockStorage реализует интерфейс repositories.Storage для тестирования
type MockStorage struct {
	shouldFail bool
	items      []models.EncryptedItem
	counts     map[models.ItemType]int32
}

func (m *MockStorage) SignUpUser(ctx context.Context, user *models.User) error { return nil }
func (m *MockStorage) GetUser(ctx context.Context, login string) (*models.User, error) {
	return nil, nil
}

func (m *MockStorage) GetAllUserItems(ctx context.Context, login string) ([]models.EncryptedItem, error) {
	if m.shouldFail {
		return nil, errors.New("storage error")
	}
	return m.items, nil
}

func (m *MockStorage) GetUserItemsWithType(ctx context.Context, typ models.ItemType, login string) ([]models.EncryptedItem, error) {
	if m.shouldFail {
		return nil, errors.New("storage error")
	}
	var filtered []models.EncryptedItem
	for _, item := range m.items {
		if item.Type == typ {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (m *MockStorage) AddItem(ctx context.Context, item *models.EncryptedItem) error {
	if m.shouldFail {
		return errors.New("storage error")
	}
	return nil
}

func (m *MockStorage) EditItem(ctx context.Context, item *models.EncryptedItem) error {
	if m.shouldFail {
		return errors.New("storage error")
	}
	return nil
}

func (m *MockStorage) DeleteItem(ctx context.Context, login string, itemID [16]byte) error {
	if m.shouldFail {
		return errors.New("storage error")
	}
	return nil
}

func (m *MockStorage) GetTypesCounts(ctx context.Context, login string) (map[models.ItemType]int32, error) {
	if m.shouldFail {
		return nil, errors.New("storage error")
	}
	return m.counts, nil
}

func TestNewItemService(t *testing.T) {
	repo := &MockStorage{}
	service := NewItemService(repo)

	assert.NotNil(t, service)
}

func TestItemService_GetUserItems(t *testing.T) {
	tests := []struct {
		name     string
		typ      models.ItemType
		login    string
		mockData []models.EncryptedItem
		wantErr  bool
	}{
		{
			name:  "get all items",
			typ:   models.ItemTypeUNSPECIFIED,
			login: "testuser",
			mockData: []models.EncryptedItem{
				{Name: "item1", Type: models.ItemTypeCREDENTIALS},
				{Name: "item2", Type: models.ItemTypeTEXT},
			},
			wantErr: false,
		},
		{
			name:  "get items with specific type",
			typ:   models.ItemTypeCREDENTIALS,
			login: "testuser",
			mockData: []models.EncryptedItem{
				{Name: "item1", Type: models.ItemTypeCREDENTIALS},
				{Name: "item2", Type: models.ItemTypeTEXT},
			},
			wantErr: false,
		},
		{
			name:     "storage error",
			typ:      models.ItemTypeUNSPECIFIED,
			login:    "testuser",
			mockData: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStorage{
				shouldFail: tt.wantErr,
				items:      tt.mockData,
			}
			service := NewItemService(mockRepo)

			items, err := service.GetUserItems(context.Background(), tt.typ, tt.login)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, items)
			}
		})
	}
}

func TestItemService_GetTypesCounts(t *testing.T) {
	tests := []struct {
		name       string
		login      string
		mockCounts map[models.ItemType]int32
		wantErr    bool
	}{
		{
			name:  "successful get counts",
			login: "testuser",
			mockCounts: map[models.ItemType]int32{
				models.ItemTypeCREDENTIALS: 5,
				models.ItemTypeTEXT:        3,
			},
			wantErr: false,
		},
		{
			name:       "storage error",
			login:      "testuser",
			mockCounts: nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStorage{
				shouldFail: tt.wantErr,
				counts:     tt.mockCounts,
			}
			service := NewItemService(mockRepo)

			counts, err := service.GetTypesCounts(context.Background(), tt.login)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, counts)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, counts)
			}
		})
	}
}

func TestItemService_AddItem(t *testing.T) {
	tests := []struct {
		name    string
		item    *models.EncryptedItem
		wantErr bool
	}{
		{
			name: "successful add",
			item: &models.EncryptedItem{
				Name: "test item",
				Type: models.ItemTypeCREDENTIALS,
			},
			wantErr: false,
		},
		{
			name: "storage error",
			item: &models.EncryptedItem{
				Name: "test item",
				Type: models.ItemTypeCREDENTIALS,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStorage{shouldFail: tt.wantErr}
			service := NewItemService(mockRepo)

			err := service.AddItem(context.Background(), tt.item)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestItemService_EditItem(t *testing.T) {
	tests := []struct {
		name    string
		item    *models.EncryptedItem
		wantErr bool
	}{
		{
			name: "successful edit",
			item: &models.EncryptedItem{
				Name: "edited item",
				Type: models.ItemTypeCREDENTIALS,
			},
			wantErr: false,
		},
		{
			name: "storage error",
			item: &models.EncryptedItem{
				Name: "edited item",
				Type: models.ItemTypeCREDENTIALS,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStorage{shouldFail: tt.wantErr}
			service := NewItemService(mockRepo)

			err := service.EditItem(context.Background(), tt.item)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestItemService_DeleteItem(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		itemID  [16]byte
		wantErr bool
	}{
		{
			name:    "successful delete",
			login:   "testuser",
			itemID:  [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			wantErr: false,
		},
		{
			name:    "storage error",
			login:   "testuser",
			itemID:  [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockStorage{shouldFail: tt.wantErr}
			service := NewItemService(mockRepo)

			err := service.DeleteItem(context.Background(), tt.login, tt.itemID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
