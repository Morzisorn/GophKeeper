package ui

import (
	"gophkeeper/models"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_handleViewItemDetails_ValidItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{
					ID:   [16]byte{1, 2, 3},
					Name: "test-item",
				},
			},
		},
		input: "some-input",
	}

	result, cmd := ui.handleViewItemDetails()

	assert.Equal(t, ui, result)
	assert.NotNil(t, cmd) // Should return decryptItemCmd
	assert.Equal(t, stateItemDetails, ui.state)
	assert.Empty(t, ui.input)
	assert.Nil(t, ui.itemCtrl.decryptedItem)
	assert.NotNil(t, ui.itemCtrl.selectedItem)
	assert.Equal(t, "test-item", ui.itemCtrl.selectedItem.Name)
}

func TestUIController_handleViewItemDetails_InvalidIndex(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 5, // Index out of bounds
			items: []models.EncryptedItem{
				{Name: "test-item"},
			},
		},
	}

	result, cmd := ui.handleViewItemDetails()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
}

func TestUIController_handleViewItemDetails_EmptyItems(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items:       []models.EncryptedItem{},
		},
	}

	result, cmd := ui.handleViewItemDetails()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
}

func TestUIController_handleItemDetailsInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleItemDetailsInput_Back_Escape(t *testing.T) {
	ui := &UIController{
		state: stateItemDetails,
	}

	model, cmd := ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemsList, ui.state)
}

func TestUIController_handleItemDetailsInput_Back_B(t *testing.T) {
	ui := &UIController{
		state: stateItemDetails,
	}

	model, cmd := ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemsList, ui.state)
}

func TestUIController_handleItemDetailsInput_Delete(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1, // Set to non-zero initially
	}

	model, cmd := ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateConfirmDelete, ui.state)
	assert.Equal(t, 0, ui.confirmChoice) // Should be reset to 0
}

func TestUIController_handleItemDetailsInput_Edit_WithDecryptedItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Name: "test-item",
				Type: models.ItemTypeTEXT,
				Data: &models.Text{Content: "test"},
			},
		},
	}

	model, cmd := ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)                        // startEditItem returns nil command
	assert.NotNil(t, ui.itemCtrl.editingItem) // Should create editing item
	assert.Equal(t, stateEditItemName, ui.state)
}

func TestUIController_handleItemDetailsInput_Edit_WithoutDecryptedItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: nil,
		},
		state: stateItemDetails,
	}

	model, cmd := ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state) // Should remain unchanged
}

func TestUIController_handleItemDetailsInput_Metadata_WithDecryptedItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Name: "test-item",
			},
		},
	}

	model, cmd := ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd) // startManageMetadata returns nil command
}

func TestUIController_handleItemDetailsInput_Metadata_WithoutDecryptedItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: nil,
		},
		state: stateItemDetails,
	}

	model, cmd := ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state) // Should remain unchanged
}

func TestUIController_handleItemDetailsInput_OtherKey(t *testing.T) {
	ui := &UIController{
		state: stateItemDetails,
	}

	model, cmd := ui.handleItemDetailsInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state) // Should remain unchanged
}

func TestUIController_itemDetailsView_NoItems(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items:       []models.EncryptedItem{},
		},
	}

	view := ui.itemDetailsView()

	assert.Equal(t, "No item selected", view)
}

func TestUIController_itemDetailsView_InvalidIndex(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 5, // Index out of bounds
			items: []models.EncryptedItem{
				{Name: "test-item"},
			},
		},
	}

	view := ui.itemDetailsView()

	assert.Equal(t, "No item selected", view)
}

func TestUIController_itemDetailsView_NoDecryptedData(t *testing.T) {
	now := time.Now()
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{
					Name:      "test-item",
					Type:      models.ItemTypeTEXT,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			decryptedItem: nil,
		},
	}

	view := ui.itemDetailsView()

	assert.Contains(t, view, "Item Details: test-item")
	assert.Contains(t, view, "Type: TEXT")
	assert.Contains(t, view, "Created:")
	assert.Contains(t, view, "Updated:")
	assert.Contains(t, view, "Loading data...")
	assert.Contains(t, view, "e to edit")
	assert.Contains(t, view, "m to manage metadata")
	assert.Contains(t, view, "d to delete")
}

func TestUIController_itemDetailsView_WithCredentials(t *testing.T) {
	now := time.Now()
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{
					Name:      "test-creds",
					Type:      models.ItemTypeCREDENTIALS,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			decryptedItem: &models.Item{
				Data: &models.Credentials{
					Login:    "test-login",
					Password: "test-password",
				},
				Meta: models.Meta{
					Map: map[string]string{
						"note": "test note",
					},
				},
			},
		},
	}

	view := ui.itemDetailsView()

	assert.Contains(t, view, "Item Details: test-creds")
	assert.Contains(t, view, "Type: CREDENTIALS")
	assert.Contains(t, view, "Data:")
	assert.Contains(t, view, "Login: test-login")
	assert.Contains(t, view, "Password: test-password")
	assert.Contains(t, view, "Metadata:")
	assert.Contains(t, view, "note: test note")
}

func TestUIController_itemDetailsView_WithText(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{
					Name: "test-text",
					Type: models.ItemTypeTEXT,
				},
			},
			decryptedItem: &models.Item{
				Data: &models.Text{
					Content: "test content",
				},
				Meta: models.Meta{
					Map: map[string]string{},
				},
			},
		},
	}

	view := ui.itemDetailsView()

	assert.Contains(t, view, "Item Details: test-text")
	assert.Contains(t, view, "Type: TEXT")
	assert.Contains(t, view, "Content: test content")
	assert.Contains(t, view, "No metadata")
}

func TestUIController_itemDetailsView_WithCard(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{
					Name: "test-card",
					Type: models.ItemTypeCARD,
				},
			},
			decryptedItem: &models.Item{
				Data: &models.Card{
					Number:         "1234567890123456",
					ExpiryDate:     "12/25",
					SecurityCode:   "123",
					CardholderName: "John Doe",
				},
				Meta: models.Meta{
					Map: map[string]string{},
				},
			},
		},
	}

	view := ui.itemDetailsView()

	assert.Contains(t, view, "Item Details: test-card")
	assert.Contains(t, view, "Type: CARD")
	assert.Contains(t, view, "Number: 1234567890123456")
	assert.Contains(t, view, "Expiry: 12/25")
	assert.Contains(t, view, "CVV: 123")
	assert.Contains(t, view, "Cardholder: John Doe")
}

func TestUIController_itemDetailsView_WithBinary(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{
					Name: "test-binary",
					Type: models.ItemTypeBINARY,
				},
			},
			decryptedItem: &models.Item{
				Data: &models.Binary{
					Content: []byte("binary data"),
				},
				Meta: models.Meta{
					Map: map[string]string{},
				},
			},
		},
	}

	view := ui.itemDetailsView()

	assert.Contains(t, view, "Item Details: test-binary")
	assert.Contains(t, view, "Type: BINARY")
	assert.Contains(t, view, "Content: binary data")
}

type UnknownItemType struct {}

func (u UnknownItemType) GetType() models.ItemType {
	return "UNKNOWN"
}

func TestUIController_itemDetailsView_UnknownDataType(t *testing.T) {
	var unknownType UnknownItemType
	ui := &UIController{
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{
					Name: "test-unknown",
					Type: models.ItemType("UNKNOWN"),
				},
			},
			decryptedItem: &models.Item{
				Data: unknownType, // Not a supported type
				Meta: models.Meta{
					Map: map[string]string{},
				},
			},
		},
	}

	view := ui.itemDetailsView()

	assert.Contains(t, view, "Item Details: test-unknown")
	assert.Contains(t, view, "Type: UNKNOWN")
	assert.Contains(t, view, "Unknown data type")
}
