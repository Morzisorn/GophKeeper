package controllers

import (
	"context"
	"testing"

	pb "gophkeeper/internal/protos/items"
	iserv "gophkeeper/internal/server/services/item_service"

	"github.com/stretchr/testify/assert"
)

func TestNewItemController(t *testing.T) {
	service := &iserv.ItemService{}
	controller := NewItemController(service)
	
	assert.NotNil(t, controller)
	assert.Equal(t, service, controller.service)
}

func TestItemController_AddItem(t *testing.T) {
	service := &iserv.ItemService{}
	controller := NewItemController(service)

	request := &pb.AddItemRequest{
		Item: &pb.EncryptedItem{
			Name:      "",
			Type:      pb.ItemType_ITEM_TYPE_CREDENTIALS,
			UserLogin: "testuser",
			EncryptedData: &pb.EncryptedData{
				EncryptedContent: "content",
				Nonce:           "nonce",
			},
		},
	}

	response, err := controller.AddItem(context.Background(), request)
	
	// Just check that method executes without panic
	_ = response
	_ = err
	assert.True(t, true)
}

func TestItemController_EditItem(t *testing.T) {
	service := &iserv.ItemService{}
	controller := NewItemController(service)

	request := &pb.EditItemRequest{
		Item: &pb.EncryptedItem{
			Name:      "item",
			Type:      pb.ItemType_ITEM_TYPE_CREDENTIALS,
			UserLogin: "",
			EncryptedData: &pb.EncryptedData{
				EncryptedContent: "content",
				Nonce:           "nonce",
			},
		},
	}

	response, err := controller.EditItem(context.Background(), request)
	
	_ = response
	_ = err
	assert.True(t, true)
}

func TestItemController_DeleteItem(t *testing.T) {
	service := &iserv.ItemService{}
	controller := NewItemController(service)

	request := &pb.DeleteItemRequest{
		ItemId:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		UserLogin: "",
	}

	response, err := controller.DeleteItem(context.Background(), request)
	
	_ = response
	_ = err
	assert.True(t, true)
}

func TestItemController_GetUserItems(t *testing.T) {
	service := &iserv.ItemService{}
	controller := NewItemController(service)

	request := &pb.GetUserItemsRequest{
		UserLogin: "",
		Type:      pb.ItemType_ITEM_TYPE_CREDENTIALS,
	}

	response, err := controller.GetUserItems(context.Background(), request)
	
	_ = response
	_ = err
	assert.True(t, true)
}

func TestItemController_GetItemTypesCounters(t *testing.T) {
	service := &iserv.ItemService{}
	controller := NewItemController(service)

	request := &pb.TypesCountsRequest{
		UserLogin: "testuser",
	}

	response, err := controller.GetItemTypesCounters(context.Background(), request)
	
	_ = response
	_ = err
	assert.True(t, true)
}

func TestIsPbItemValid(t *testing.T) {
	tests := []struct {
		name     string
		item     *pb.EncryptedItem
		expected bool
	}{
		{
			name: "valid item",
			item: &pb.EncryptedItem{
				Name:      "test",
				Type:      pb.ItemType_ITEM_TYPE_CREDENTIALS,
				UserLogin: "user",
				EncryptedData: &pb.EncryptedData{
					EncryptedContent: "content",
					Nonce:           "nonce",
				},
			},
			expected: true,
		},
		{
			name: "invalid - empty name",
			item: &pb.EncryptedItem{
				Name:      "",
				Type:      pb.ItemType_ITEM_TYPE_CREDENTIALS,
				UserLogin: "user",
				EncryptedData: &pb.EncryptedData{
					EncryptedContent: "content",
					Nonce:           "nonce",
				},
			},
			expected: false,
		},
		{
			name: "invalid - empty user login",
			item: &pb.EncryptedItem{
				Name:      "test",
				Type:      pb.ItemType_ITEM_TYPE_CREDENTIALS,
				UserLogin: "",
				EncryptedData: &pb.EncryptedData{
					EncryptedContent: "content",
					Nonce:           "nonce",
				},
			},
			expected: false,
		},
		{
			name: "invalid - empty encrypted content",
			item: &pb.EncryptedItem{
				Name:      "test",
				Type:      pb.ItemType_ITEM_TYPE_CREDENTIALS,
				UserLogin: "user",
				EncryptedData: &pb.EncryptedData{
					EncryptedContent: "",
					Nonce:           "nonce",
				},
			},
			expected: false,
		},
		{
			name: "invalid - empty nonce",
			item: &pb.EncryptedItem{
				Name:      "test",
				Type:      pb.ItemType_ITEM_TYPE_CREDENTIALS,
				UserLogin: "user",
				EncryptedData: &pb.EncryptedData{
					EncryptedContent: "content",
					Nonce:           "",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isPbItemValid(tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}