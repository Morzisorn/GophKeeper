package models

import (
	"testing"
	"time"

	pb "gophkeeper/internal/protos/items"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEncryptedItemPbToModels(t *testing.T) {
	now := time.Now()
	pbItem := &pb.EncryptedItem{
		Id:        []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		UserLogin: "testuser",
		Name:      "test item",
		Type:      pb.ItemType_ITEM_TYPE_CREDENTIALS,
		EncryptedData: &pb.EncryptedData{
			EncryptedContent: "encrypted_content",
			Nonce:           "test_nonce",
		},
		Meta:      map[string]string{"key": "value"},
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}

	result := EncryptedItemPbToModels(pbItem)

	require.NotNil(t, result)
	assert.Equal(t, [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, result.ID)
	assert.Equal(t, "testuser", result.UserLogin)
	assert.Equal(t, "test item", result.Name)
	assert.Equal(t, ItemTypeCREDENTIALS, result.Type)
	assert.Equal(t, "encrypted_content", result.EncryptedData.EncryptedContent)
	assert.Equal(t, "test_nonce", result.EncryptedData.Nonce)
	assert.Equal(t, map[string]string{"key": "value"}, result.Meta.Map)
	assert.Equal(t, now.Unix(), result.CreatedAt.Unix())
	assert.Equal(t, now.Unix(), result.UpdatedAt.Unix())
}

func TestItemIdPbToModels(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected [16]byte
	}{
		{
			name:     "full 16 bytes",
			input:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			expected: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
		{
			name:     "less than 16 bytes",
			input:    []byte{1, 2, 3},
			expected: [16]byte{1, 2, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "more than 16 bytes",
			input:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18},
			expected: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ItemIdPbToModels(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEncryptedItem_ToPb(t *testing.T) {
	now := time.Now()
	item := &EncryptedItem{
		ID:        [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		UserLogin: "testuser",
		Name:      "test item",
		Type:      ItemTypeCREDENTIALS,
		EncryptedData: EncryptedData{
			EncryptedContent: "encrypted_content",
			Nonce:           "test_nonce",
		},
		Meta:      Meta{Map: map[string]string{"key": "value"}},
		CreatedAt: now,
		UpdatedAt: now,
	}

	result, err := item.ToPb()

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, result.Id)
	assert.Equal(t, "testuser", result.UserLogin)
	assert.Equal(t, "test item", result.Name)
	assert.Equal(t, pb.ItemType_ITEM_TYPE_CREDENTIALS, result.Type)
	assert.Equal(t, "encrypted_content", result.EncryptedData.EncryptedContent)
	assert.Equal(t, "test_nonce", result.EncryptedData.Nonce)
	assert.Equal(t, map[string]string{"key": "value"}, result.Meta)
	assert.Equal(t, now.Unix(), result.CreatedAt.AsTime().Unix())
	assert.Equal(t, now.Unix(), result.UpdatedAt.AsTime().Unix())
}

func TestEncryptedData_ToPb(t *testing.T) {
	ed := &EncryptedData{
		EncryptedContent: "test_content",
		Nonce:           "test_nonce",
	}

	result := ed.ToPb()

	require.NotNil(t, result)
	assert.Equal(t, "test_content", result.EncryptedContent)
	assert.Equal(t, "test_nonce", result.Nonce)
}

func TestEncryptedDataPbToModel(t *testing.T) {
	pbData := &pb.EncryptedData{
		EncryptedContent: "test_content",
		Nonce:           "test_nonce",
	}

	result := EncryptedDataPbToModel(pbData)

	assert.Equal(t, "test_content", result.EncryptedContent)
	assert.Equal(t, "test_nonce", result.Nonce)
}

func TestItemTypePbToModel(t *testing.T) {
	tests := []struct {
		name     string
		input    pb.ItemType
		expected ItemType
	}{
		{
			name:     "unspecified",
			input:    pb.ItemType_ITEM_TYPE_UNSPECIFIED,
			expected: ItemTypeUNSPECIFIED,
		},
		{
			name:     "credentials",
			input:    pb.ItemType_ITEM_TYPE_CREDENTIALS,
			expected: ItemTypeCREDENTIALS,
		},
		{
			name:     "text",
			input:    pb.ItemType_ITEM_TYPE_TEXT,
			expected: ItemTypeTEXT,
		},
		{
			name:     "binary",
			input:    pb.ItemType_ITEM_TYPE_BINARY,
			expected: ItemTypeBINARY,
		},
		{
			name:     "card",
			input:    pb.ItemType_ITEM_TYPE_CARD,
			expected: ItemTypeCARD,
		},
		{
			name:     "unknown value",
			input:    pb.ItemType(999),
			expected: ItemTypeUNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ItemTypePbToModel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestItemType_ToPb(t *testing.T) {
	tests := []struct {
		name     string
		input    ItemType
		expected pb.ItemType
	}{
		{
			name:     "unspecified",
			input:    ItemTypeUNSPECIFIED,
			expected: pb.ItemType_ITEM_TYPE_UNSPECIFIED,
		},
		{
			name:     "credentials",
			input:    ItemTypeCREDENTIALS,
			expected: pb.ItemType_ITEM_TYPE_CREDENTIALS,
		},
		{
			name:     "text",
			input:    ItemTypeTEXT,
			expected: pb.ItemType_ITEM_TYPE_TEXT,
		},
		{
			name:     "binary",
			input:    ItemTypeBINARY,
			expected: pb.ItemType_ITEM_TYPE_BINARY,
		},
		{
			name:     "card",
			input:    ItemTypeCARD,
			expected: pb.ItemType_ITEM_TYPE_CARD,
		},
		{
			name:     "unknown value",
			input:    ItemType("UNKNOWN"),
			expected: pb.ItemType_ITEM_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToPb()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	// Тест полного цикла: models -> pb -> models
	original := &EncryptedItem{
		ID:        [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		UserLogin: "testuser",
		Name:      "test item",
		Type:      ItemTypeCREDENTIALS,
		EncryptedData: EncryptedData{
			EncryptedContent: "encrypted_content",
			Nonce:           "test_nonce",
		},
		Meta:      Meta{Map: map[string]string{"key": "value"}},
		CreatedAt: time.Now().Truncate(time.Second), // Обрезаем до секунд для сравнения
		UpdatedAt: time.Now().Truncate(time.Second),
	}

	// models -> pb
	pbItem, err := original.ToPb()
	require.NoError(t, err)
	require.NotNil(t, pbItem)

	// pb -> models
	converted := EncryptedItemPbToModels(pbItem)
	require.NotNil(t, converted)

	// Сравниваем
	assert.Equal(t, original.ID, converted.ID)
	assert.Equal(t, original.UserLogin, converted.UserLogin)
	assert.Equal(t, original.Name, converted.Name)
	assert.Equal(t, original.Type, converted.Type)
	assert.Equal(t, original.EncryptedData, converted.EncryptedData)
	assert.Equal(t, original.Meta.Map, converted.Meta.Map)
	assert.Equal(t, original.CreatedAt.Unix(), converted.CreatedAt.Unix())
	assert.Equal(t, original.UpdatedAt.Unix(), converted.UpdatedAt.Unix())
}