package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemType_String(t *testing.T) {
	tests := []struct {
		name     string
		itemType ItemType
		expected string
	}{
		{
			name:     "credentials type",
			itemType: ItemTypeCREDENTIALS,
			expected: "CREDENTIALS",
		},
		{
			name:     "text type",
			itemType: ItemTypeTEXT,
			expected: "TEXT",
		},
		{
			name:     "binary type",
			itemType: ItemTypeBINARY,
			expected: "BINARY",
		},
		{
			name:     "card type",
			itemType: ItemTypeCARD,
			expected: "CARD",
		},
		{
			name:     "unspecified type",
			itemType: ItemTypeUNSPECIFIED,
			expected: "UNSPECIFIED",
		},
		{
			name:     "custom type",
			itemType: ItemType("CUSTOM"),
			expected: "CUSTOM",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.itemType.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCredentials_GetType(t *testing.T) {
	creds := Credentials{
		Login:    "testuser",
		Password: "testpass",
	}

	result := creds.GetType()
	assert.Equal(t, ItemTypeCREDENTIALS, result)
}

func TestText_GetType(t *testing.T) {
	text := Text{
		Content: "test content",
	}

	result := text.GetType()
	assert.Equal(t, ItemTypeTEXT, result)
}

func TestBinary_GetType(t *testing.T) {
	binary := Binary{
		Content: []byte("test binary data"),
	}

	result := binary.GetType()
	assert.Equal(t, ItemTypeBINARY, result)
}

func TestCard_GetType(t *testing.T) {
	card := Card{
		Number:         "1234567890123456",
		ExpiryDate:     "12/25",
		SecurityCode:   "123",
		CardholderName: "Test User",
	}

	result := card.GetType()
	assert.Equal(t, ItemTypeCARD, result)
}

func TestItemType_CreateDataByType(t *testing.T) {
	tests := []struct {
		name         string
		itemType     ItemType
		expectedType ItemType
		wantErr      bool
	}{
		{
			name:         "create credentials",
			itemType:     ItemTypeCREDENTIALS,
			expectedType: ItemTypeCREDENTIALS,
			wantErr:      false,
		},
		{
			name:         "create text",
			itemType:     ItemTypeTEXT,
			expectedType: ItemTypeTEXT,
			wantErr:      false,
		},
		{
			name:         "create binary",
			itemType:     ItemTypeBINARY,
			expectedType: ItemTypeBINARY,
			wantErr:      false,
		},
		{
			name:         "create card",
			itemType:     ItemTypeCARD,
			expectedType: ItemTypeCARD,
			wantErr:      false,
		},
		{
			name:     "unknown type",
			itemType: ItemType("UNKNOWN"),
			wantErr:  true,
		},
		{
			name:     "unspecified type",
			itemType: ItemTypeUNSPECIFIED,
			wantErr:  true,
		},
		{
			name:     "empty type",
			itemType: ItemType(""),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.itemType.CreateDataByType()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, data)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, data)
				assert.Equal(t, tt.expectedType, data.GetType())
			}
		})
	}
}

func TestItemType_CreateDataByType_Types(t *testing.T) {
	// Тестируем что возвращаются правильные типы структур
	t.Run("credentials returns Credentials struct", func(t *testing.T) {
		data, err := ItemTypeCREDENTIALS.CreateDataByType()
		require.NoError(t, err)
		_, ok := data.(*Credentials)
		assert.True(t, ok)
	})

	t.Run("text returns Text struct", func(t *testing.T) {
		data, err := ItemTypeTEXT.CreateDataByType()
		require.NoError(t, err)
		_, ok := data.(*Text)
		assert.True(t, ok)
	})

	t.Run("binary returns Binary struct", func(t *testing.T) {
		data, err := ItemTypeBINARY.CreateDataByType()
		require.NoError(t, err)
		_, ok := data.(*Binary)
		assert.True(t, ok)
	})

	t.Run("card returns Card struct", func(t *testing.T) {
		data, err := ItemTypeCARD.CreateDataByType()
		require.NoError(t, err)
		_, ok := data.(*Card)
		assert.True(t, ok)
	})
}

func TestItemTypes_Constant(t *testing.T) {
	// Тестируем что константа ItemTypes содержит все ожидаемые типы
	expectedTypes := []ItemType{
		ItemTypeCREDENTIALS,
		ItemTypeTEXT,
		ItemTypeBINARY,
		ItemTypeCARD,
	}

	assert.Len(t, ItemTypes, len(expectedTypes))

	for _, expectedType := range expectedTypes {
		assert.Contains(t, ItemTypes, expectedType)
	}

	// Проверяем что UNSPECIFIED не включен в ItemTypes
	assert.NotContains(t, ItemTypes, ItemTypeUNSPECIFIED)
}

func TestData_Interface(t *testing.T) {
	// Тестируем что все структуры данных реализуют интерфейс Data
	var data Data

	data = &Credentials{}
	assert.Equal(t, ItemTypeCREDENTIALS, data.GetType())

	data = &Text{}
	assert.Equal(t, ItemTypeTEXT, data.GetType())

	data = &Binary{}
	assert.Equal(t, ItemTypeBINARY, data.GetType())

	data = &Card{}
	assert.Equal(t, ItemTypeCARD, data.GetType())
}

func TestItemType_Constants(t *testing.T) {
	// Тестируем что константы имеют правильные значения
	assert.Equal(t, "UNSPECIFIED", string(ItemTypeUNSPECIFIED))
	assert.Equal(t, "CREDENTIALS", string(ItemTypeCREDENTIALS))
	assert.Equal(t, "TEXT", string(ItemTypeTEXT))
	assert.Equal(t, "BINARY", string(ItemTypeBINARY))
	assert.Equal(t, "CARD", string(ItemTypeCARD))
}

func TestDataStructures_DefaultValues(t *testing.T) {
	// Тестируем что структуры данных могут быть созданы с нулевыми значениями
	t.Run("credentials with default values", func(t *testing.T) {
		creds := Credentials{}
		assert.Equal(t, "", creds.Login)
		assert.Equal(t, "", creds.Password)
		assert.Equal(t, ItemTypeCREDENTIALS, creds.GetType())
	})

	t.Run("text with default values", func(t *testing.T) {
		text := Text{}
		assert.Equal(t, "", text.Content)
		assert.Equal(t, ItemTypeTEXT, text.GetType())
	})

	t.Run("binary with default values", func(t *testing.T) {
		binary := Binary{}
		assert.Nil(t, binary.Content)
		assert.Equal(t, ItemTypeBINARY, binary.GetType())
	})

	t.Run("card with default values", func(t *testing.T) {
		card := Card{}
		assert.Equal(t, "", card.Number)
		assert.Equal(t, "", card.ExpiryDate)
		assert.Equal(t, "", card.SecurityCode)
		assert.Equal(t, "", card.CardholderName)
		assert.Equal(t, ItemTypeCARD, card.GetType())
	})
}
