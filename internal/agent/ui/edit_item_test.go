package ui

import (
	"gophkeeper/models"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_startEditItem_NilDecryptedItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: nil,
		},
	}

	result, cmd := ui.startEditItem()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
}

func TestUIController_startEditItem_Credentials(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				ID:        [16]byte{1, 2, 3},
				UserLogin: "test-user",
				Name:      "test-creds",
				Type:      models.ItemTypeCREDENTIALS,
				Data: &models.Credentials{
					Login:    "test-login",
					Password: "test-password",
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
		input: "some-input",
	}

	result, cmd := ui.startEditItem()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.NotNil(t, ui.itemCtrl.editingItem)
	assert.Equal(t, "test-creds", ui.itemCtrl.editingItem.Name)
	assert.Equal(t, models.ItemTypeCREDENTIALS, ui.itemCtrl.editingItem.Type)
	assert.IsType(t, &models.Credentials{}, ui.itemCtrl.editingItem.Data)
	assert.Equal(t, "test-login", ui.itemCtrl.editingItem.Data.(*models.Credentials).Login)
	assert.Equal(t, stateEditItemName, ui.state)
	assert.Equal(t, "test-creds", ui.input)
	assert.Equal(t, 0, ui.itemCtrl.editStep)
}

func TestUIController_startEditItem_Text(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Name: "test-text",
				Type: models.ItemTypeTEXT,
				Data: &models.Text{
					Content: "test content",
				},
			},
		},
	}

	result, cmd := ui.startEditItem()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.NotNil(t, ui.itemCtrl.editingItem)
	assert.Equal(t, models.ItemTypeTEXT, ui.itemCtrl.editingItem.Type)
	assert.IsType(t, &models.Text{}, ui.itemCtrl.editingItem.Data)
	assert.Equal(t, "test content", ui.itemCtrl.editingItem.Data.(*models.Text).Content)
}

func TestUIController_startEditItem_Card(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Name: "test-card",
				Type: models.ItemTypeCARD,
				Data: &models.Card{
					Number:         "1234567890123456",
					ExpiryDate:     "12/25",
					SecurityCode:   "123",
					CardholderName: "John Doe",
				},
			},
		},
	}

	result, cmd := ui.startEditItem()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.NotNil(t, ui.itemCtrl.editingItem)
	assert.Equal(t, models.ItemTypeCARD, ui.itemCtrl.editingItem.Type)
	assert.IsType(t, &models.Card{}, ui.itemCtrl.editingItem.Data)
	cardData := ui.itemCtrl.editingItem.Data.(*models.Card)
	assert.Equal(t, "1234567890123456", cardData.Number)
	assert.Equal(t, "12/25", cardData.ExpiryDate)
}

func TestUIController_startEditItem_Binary(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptedItem: &models.Item{
				Name: "test-binary",
				Type: models.ItemTypeBINARY,
				Data: &models.Binary{
					Content: []byte("binary data"),
				},
			},
		},
	}

	result, cmd := ui.startEditItem()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.NotNil(t, ui.itemCtrl.editingItem)
	assert.Equal(t, models.ItemTypeBINARY, ui.itemCtrl.editingItem.Type)
	assert.IsType(t, &models.Binary{}, ui.itemCtrl.editingItem.Data)
	assert.Equal(t, []byte("binary data"), ui.itemCtrl.editingItem.Data.(*models.Binary).Content)
}

func TestUIController_handleEditItemNameInput_Quit(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	model, cmd = ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleEditItemNameInput_Escape(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			editingItem: &models.Item{Name: "test"},
		},
		input: "test-input",
		state: stateEditItemName,
	}

	model, cmd := ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Nil(t, ui.itemCtrl.editingItem)
	assert.Equal(t, stateItemDetails, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_handleEditItemNameInput_Enter_ValidName(t *testing.T) {
	ui := &UIController{
		input: "  new-name  ",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Type: models.ItemTypeCREDENTIALS,
				Data: &models.Credentials{Login: "test-login"},
			},
		},
	}

	model, cmd := ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "new-name", ui.itemCtrl.editingItem.Name)
	assert.Equal(t, "test-login", ui.input) // Input должен быть установлен в значение логина
	assert.Equal(t, stateEditCredentialLogin, ui.state)
}

func TestUIController_handleEditItemNameInput_Enter_EmptyName(t *testing.T) {
	ui := &UIController{
		input: "   ",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{Name: "original"},
		},
		state: stateEditItemName,
	}

	model, cmd := ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "original", ui.itemCtrl.editingItem.Name) // Should remain unchanged
	assert.Equal(t, stateEditItemName, ui.state)              // Should remain unchanged
}

func TestUIController_handleEditItemNameInput_Enter_TextType(t *testing.T) {
	ui := &UIController{
		input: "new-text-name",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Type: models.ItemTypeTEXT,
				Data: &models.Text{Content: "test content"},
			},
		},
	}

	model, cmd := ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "new-text-name", ui.itemCtrl.editingItem.Name)
	assert.Equal(t, stateEditTextContent, ui.state)
	assert.Equal(t, "test content", ui.input)
}

func TestUIController_handleEditItemNameInput_Enter_CardType(t *testing.T) {
	ui := &UIController{
		input: "new-card-name",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Type: models.ItemTypeCARD,
				Data: &models.Card{Number: "1234567890123456"},
			},
		},
	}

	model, cmd := ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "new-card-name", ui.itemCtrl.editingItem.Name)
	assert.Equal(t, stateEditCardNumber, ui.state)
	assert.Equal(t, "1234567890123456", ui.input)
}

func TestUIController_handleEditItemNameInput_Enter_BinaryType(t *testing.T) {
	ui := &UIController{
		input: "new-binary-name",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Type: models.ItemTypeBINARY,
				Data: &models.Binary{Content: []byte("binary data")},
			},
		},
	}

	model, cmd := ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "new-binary-name", ui.itemCtrl.editingItem.Name)
	assert.Equal(t, stateEditBinaryData, ui.state)
	assert.Equal(t, "binary data", ui.input)
}

func TestUIController_handleEditItemNameInput_Backspace(t *testing.T) {
	ui := &UIController{input: "test"}

	model, cmd := ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyBackspace})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "tes", ui.input)
}

func TestUIController_handleEditItemNameInput_CharacterInput(t *testing.T) {
	ui := &UIController{input: "test"}

	model, cmd := ui.handleEditItemNameInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "testa", ui.input)
}

func TestUIController_handleEditCredentialLoginInput_Enter(t *testing.T) {
	ui := &UIController{
		input: "  new-login  ",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Data: &models.Credentials{Password: "test-password"},
			},
		},
	}

	model, cmd := ui.handleEditCredentialLoginInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "new-login", ui.itemCtrl.editingItem.Data.(*models.Credentials).Login)
	assert.Equal(t, stateEditCredentialPassword, ui.state)
	assert.Equal(t, "test-password", ui.input)
}

func TestUIController_handleEditCredentialPasswordInput_Enter(t *testing.T) {
	ui := &UIController{
		input: "  new-password  ",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Data: &models.Credentials{},
			},
		},
	}

	model, cmd := ui.handleEditCredentialPasswordInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return saveEditedItemCmd
	assert.Equal(t, "new-password", ui.itemCtrl.editingItem.Data.(*models.Credentials).Password)
	assert.Empty(t, ui.input)
}

func TestUIController_saveEditedItemCmd(t *testing.T) {
	// Этот тест сложен для unit-тестирования, так как требует mock Item service
	// Пропускаем, так как это интеграционный тест
	t.Skip("Requires mock Item service - integration test")
}

func TestUIController_handleEditSuccessInput_Enter(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			editSuccessMsg: "test message",
			editingItem:    &models.Item{Name: "test"},
		},
	}

	model, cmd := ui.handleEditSuccessInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return loadItemsCmd
	assert.Equal(t, stateProcessing, ui.state)
	assert.Empty(t, ui.itemCtrl.editSuccessMsg)
	assert.Nil(t, ui.itemCtrl.editingItem)
}

func TestUIController_handleEditErrorInput_Enter(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			editErrorMsg: "test error",
			editingItem:  &models.Item{Name: "test"},
		},
	}

	model, cmd := ui.handleEditErrorInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state)
	assert.Empty(t, ui.itemCtrl.editErrorMsg)
	assert.Nil(t, ui.itemCtrl.editingItem)
}

func TestUIController_handleEditErrorInput_Escape(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			editErrorMsg: "test error",
			editingItem:  &models.Item{Name: "test"},
		},
	}

	model, cmd := ui.handleEditErrorInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state)
	assert.Empty(t, ui.itemCtrl.editErrorMsg)
	assert.Nil(t, ui.itemCtrl.editingItem)
}

func TestUIController_handleEditTextContentInput_Enter_ValidContent(t *testing.T) {
	ui := &UIController{
		input: "  new content  ",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Data: &models.Text{},
			},
		},
	}

	model, cmd := ui.handleEditTextContentInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return saveEditedItemCmd
	assert.Equal(t, "new content", ui.itemCtrl.editingItem.Data.(*models.Text).Content)
	assert.Empty(t, ui.input)
}

func TestUIController_handleEditTextContentInput_Enter_EmptyContent(t *testing.T) {
	ui := &UIController{
		input: "   ",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Data: &models.Text{Content: "original"},
			},
		},
		state: stateEditTextContent,
	}

	model, cmd := ui.handleEditTextContentInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "original", ui.itemCtrl.editingItem.Data.(*models.Text).Content) // Should remain unchanged
	assert.Equal(t, stateEditTextContent, ui.state)                                  // Should remain unchanged
}

func TestUIController_handleEditBinaryDataInput_Enter_ValidData(t *testing.T) {
	ui := &UIController{
		input: "  new binary data  ",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Data: &models.Binary{},
			},
		},
	}

	model, cmd := ui.handleEditBinaryDataInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return saveEditedItemCmd
	assert.Equal(t, []byte("new binary data"), ui.itemCtrl.editingItem.Data.(*models.Binary).Content)
	assert.Empty(t, ui.input)
}

// View tests
func TestUIController_editItemNameView(t *testing.T) {
	ui := &UIController{
		input: "test-name",
	}

	view := ui.editItemNameView()

	assert.Contains(t, view, "Edit Item - Name")
	assert.Contains(t, view, "test-name")
	assert.Contains(t, view, "█")
	assert.Contains(t, view, "Esc to cancel")
}

func TestUIController_editCredentialLoginView(t *testing.T) {
	ui := &UIController{
		input: "test-login",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Name: "test-item",
			},
		},
	}

	view := ui.editCredentialLoginView()

	assert.Contains(t, view, "Edit Credentials - Login")
	assert.Contains(t, view, "test-item")
	assert.Contains(t, view, "test-login")
	assert.Contains(t, view, "█")
}

func TestUIController_editCredentialPasswordView(t *testing.T) {
	ui := &UIController{
		input: "password123",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Name: "test-item",
				Data: &models.Credentials{Login: "test-login"},
			},
		},
	}

	view := ui.editCredentialPasswordView()

	assert.Contains(t, view, "Edit Credentials - Password")
	assert.Contains(t, view, "test-item")
	assert.Contains(t, view, "test-login")
	assert.Contains(t, view, "***********")    // Hidden password
	assert.NotContains(t, view, "password123") // Should not show actual password
	assert.Contains(t, view, "█")
}

func TestUIController_editSuccessView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			editSuccessMsg: "Item updated successfully!",
		},
	}

	view := ui.editSuccessView()

	assert.Contains(t, view, "Item Updated Successfully")
	assert.Contains(t, view, "Item updated successfully!")
	assert.Contains(t, view, "Enter to return to items list")
}

func TestUIController_editErrorView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			editErrorMsg: "Failed to update item",
		},
	}

	view := ui.editErrorView()

	assert.Contains(t, view, "Edit Error")
	assert.Contains(t, view, "Failed to update item")
	assert.Contains(t, view, "Enter to return to item details")
}

func TestUIController_editTextContentView(t *testing.T) {
	ui := &UIController{
		input: "test content",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Name: "test-text",
			},
		},
	}

	view := ui.editTextContentView()

	assert.Contains(t, view, "Edit Text - Content")
	assert.Contains(t, view, "test-text")
	assert.Contains(t, view, "test content")
	assert.Contains(t, view, "█")
}

func TestUIController_editBinaryDataView(t *testing.T) {
	ui := &UIController{
		input: "binary data",
		itemCtrl: itemCtrl{
			editingItem: &models.Item{
				Name: "test-binary",
			},
		},
	}

	view := ui.editBinaryDataView()

	assert.Contains(t, view, "Edit Binary - Data")
	assert.Contains(t, view, "test-binary")
	assert.Contains(t, view, "binary data")
	assert.Contains(t, view, "█")
}
