package ui

import (
	"gophkeeper/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_handleAddItem(t *testing.T) {
	ui := &UIController{
		userCtrl: userCtrl{
			login: "test-user",
		},
		state: stateMenuLoggedIn,
		input: "some-input",
	}

	result, cmd := ui.handleAddItem()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, stateAddItem, ui.state)
	assert.Empty(t, ui.input)
	assert.Equal(t, 0, ui.itemTypeMenu)
	assert.Equal(t, 3, ui.maxItemTypes)
	assert.Equal(t, "test-user", ui.newItem.UserLogin)
	assert.Empty(t, ui.addItemErrorMsg)
	assert.Empty(t, ui.addItemSuccessMsg)
}

func TestUIController_handleAddItemSuccessInput_Quit(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleAddItemSuccessInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	model, cmd = ui.handleAddItemSuccessInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleAddItemSuccessInput_Enter(t *testing.T) {
	ui := &UIController{
		state: stateAddItemSuccess,
		itemCtrl: itemCtrl{
			addItemSuccessMsg: "test message",
		},
		loggedInMenu: 5,
	}

	model, cmd := ui.handleAddItemSuccessInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Empty(t, ui.addItemSuccessMsg)
	assert.Equal(t, 0, ui.loggedInMenu)
}

func TestUIController_handleAddItemErrorInput_Quit(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleAddItemErrorInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	model, cmd = ui.handleAddItemErrorInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleAddItemErrorInput_Enter(t *testing.T) {
	ui := &UIController{
		state: stateAddItemError,
		itemCtrl: itemCtrl{
			addItemErrorMsg: "test error",
			newItem:         models.Item{Name: "test"},
		},
		loggedInMenu: 5,
	}

	model, cmd := ui.handleAddItemErrorInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Empty(t, ui.addItemErrorMsg)
	assert.Equal(t, models.Item{}, ui.newItem)
	assert.Equal(t, 0, ui.loggedInMenu)
}

func TestUIController_handleItemTypeSelection_Quit(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleItemTypeSelection(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleItemTypeSelection_Escape(t *testing.T) {
	ui := &UIController{
		state: stateAddItem,
	}
	ui.messages.init()

	model, cmd := ui.handleItemTypeSelection(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
}

func TestUIController_handleItemTypeSelection_Navigation(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			itemTypeMenu: 1,
			maxItemTypes: 3,
		},
	}

	// Test up
	model, cmd := ui.handleItemTypeSelection(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.itemTypeMenu)

	// Test down
	model, cmd = ui.handleItemTypeSelection(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.itemTypeMenu)
}

func TestUIController_handleItemTypeSelection_DirectSelection(t *testing.T) {
	ui := &UIController{}

	// Test selecting option 1 (Credentials)
	model, cmd := ui.handleItemTypeSelection(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.itemTypeMenu)
	assert.Equal(t, models.ItemTypeCREDENTIALS, ui.newItem.Type)
	assert.Equal(t, stateAddItemName, ui.state)

	ui = &UIController{} // Reset

	// Test selecting option 2 (Text)
	model, cmd = ui.handleItemTypeSelection(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.itemTypeMenu)
	assert.Equal(t, models.ItemTypeTEXT, ui.newItem.Type)
	assert.Equal(t, stateAddItemName, ui.state)
}

func TestUIController_selectItemType_Credentials(t *testing.T) {
	ui := &UIController{
		input: "test-input",
	}

	result := ui.selectItemType(models.ItemTypeCREDENTIALS)

	assert.Equal(t, ui, result)
	assert.Equal(t, models.ItemTypeCREDENTIALS, ui.newItem.Type)
	assert.IsType(t, &models.Credentials{}, ui.newItem.Data)
	assert.Equal(t, stateAddItemName, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_selectItemType_Card(t *testing.T) {
	ui := &UIController{}

	result := ui.selectItemType(models.ItemTypeCARD)

	assert.Equal(t, ui, result)
	assert.Equal(t, models.ItemTypeCARD, ui.newItem.Type)
	assert.IsType(t, &models.Card{}, ui.newItem.Data)
	assert.Equal(t, stateAddItemName, ui.state)
}

func TestUIController_selectItemType_Text(t *testing.T) {
	ui := &UIController{}

	result := ui.selectItemType(models.ItemTypeTEXT)

	assert.Equal(t, ui, result)
	assert.Equal(t, models.ItemTypeTEXT, ui.newItem.Type)
	assert.IsType(t, &models.Text{}, ui.newItem.Data)
	assert.Equal(t, stateAddItemName, ui.state)
}

func TestUIController_selectItemType_Binary(t *testing.T) {
	ui := &UIController{}

	result := ui.selectItemType(models.ItemTypeBINARY)

	assert.Equal(t, ui, result)
	assert.Equal(t, models.ItemTypeBINARY, ui.newItem.Type)
	assert.IsType(t, &models.Binary{}, ui.newItem.Data)
	assert.Equal(t, stateAddItemName, ui.state)
}

func TestUIController_handleItemNameInput_Quit(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleItemNameInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleItemNameInput_Escape(t *testing.T) {
	ui := &UIController{
		state: stateAddItemName,
		input: "test-input",
	}

	model, cmd := ui.handleItemNameInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateAddItem, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_handleItemNameInput_Enter_ValidName(t *testing.T) {
	ui := &UIController{
		input: "  test-item-name  ",
		state: stateAddItemName,
	}
	ui.messages.init()

	model, cmd := ui.handleItemNameInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "test-item-name", ui.newItem.Name)
	assert.Empty(t, ui.input)
	assert.Equal(t, stateAddItemData, ui.state)
}

func TestUIController_handleItemNameInput_Enter_EmptyName(t *testing.T) {
	ui := &UIController{
		input: "   ",
		state: stateAddItemName,
	}
	ui.messages.init()

	model, cmd := ui.handleItemNameInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Empty(t, ui.newItem.Name)
	assert.Equal(t, stateAddItemName, ui.state)
}

func TestUIController_handleItemNameInput_Backspace(t *testing.T) {
	ui := &UIController{input: "test"}

	model, cmd := ui.handleItemNameInput(tea.KeyMsg{Type: tea.KeyBackspace})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "tes", ui.input)
}

func TestUIController_handleItemNameInput_CharacterInput(t *testing.T) {
	ui := &UIController{input: "test"}

	model, cmd := ui.handleItemNameInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "testa", ui.input)
}

// Validation tests
func TestValidateCardNumber_Valid(t *testing.T) {
	// Test valid 16-digit card
	err := validateCardNumber("1234567890123456")
	assert.NoError(t, err)

	// Test valid 18-digit card
	err = validateCardNumber("123456789012345678")
	assert.NoError(t, err)

	// Test with spaces and dashes
	err = validateCardNumber("1234 5678 9012 3456")
	assert.NoError(t, err)

	err = validateCardNumber("1234-5678-9012-3456")
	assert.NoError(t, err)
}

func TestValidateCardNumber_Invalid(t *testing.T) {
	// Test invalid length
	err := validateCardNumber("12345")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be 16 or 18 digits")

	// Test non-numeric
	err = validateCardNumber("abcd5678901234567")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must contain only digits")

	// Test empty
	err = validateCardNumber("")
	assert.Error(t, err)
}

func TestValidateExpiry_Valid(t *testing.T) {
	err := validateExpiry("12/25")
	assert.NoError(t, err)

	err = validateExpiry("01/30")
	assert.NoError(t, err)
}

func TestValidateExpiry_Invalid(t *testing.T) {
	// Test wrong format
	err := validateExpiry("1225")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MM/YY format")

	err = validateExpiry("12-25")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MM/YY format")

	// Test invalid month
	err = validateExpiry("13/25")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "month must be between 01 and 12")

	err = validateExpiry("00/25")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "month must be between 01 and 12")
}

func TestValidateCVV_Valid(t *testing.T) {
	err := validateCVV("123")
	assert.NoError(t, err)

	err = validateCVV("000")
	assert.NoError(t, err)
}

func TestValidateCVV_Invalid(t *testing.T) {
	// Test wrong length
	err := validateCVV("12")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exactly 3 digits")

	err = validateCVV("1234")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exactly 3 digits")

	// Test non-numeric
	err = validateCVV("12a")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exactly 3 digits")
}

// View tests
func TestUIController_addItemTypeView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			itemTypeMenu: 1,
		},
	}

	view := ui.addItemTypeView()

	assert.Contains(t, view, "Add New Item")
	assert.Contains(t, view, "Credentials")
	assert.Contains(t, view, "Text")
	assert.Contains(t, view, "Binary")
	assert.Contains(t, view, "Credit Card")
	assert.Contains(t, view, "↑/↓ to navigate")
}

func TestUIController_addItemNameView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			newItem: models.Item{Type: models.ItemTypeCREDENTIALS},
		},
		input: "test-name",
	}
	ui.messages.init()

	view := ui.addItemNameView()

	assert.Contains(t, view, "Add CREDENTIALS")
	assert.Contains(t, view, "test-name")
	assert.Contains(t, view, "█")
}

func TestUIController_addItemDataView_Credentials(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			newItem: models.Item{
				Name: "test-item",
				Type: models.ItemTypeCREDENTIALS,
			},
		},
		input: "test-login",
	}
	ui.messages.init()

	view := ui.addItemDataView()

	assert.Contains(t, view, "Add CREDENTIALS")
	assert.Contains(t, view, "test-item")
	assert.Contains(t, view, "Login")
	assert.Contains(t, view, "test-login")
}

func TestUIController_addItemDataView_Card(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			newItem: models.Item{
				Name: "test-card",
				Type: models.ItemTypeCARD,
			},
		},
		input: "1234567890123456",
	}
	ui.messages.init()

	view := ui.addItemDataView()

	assert.Contains(t, view, "Add CARD")
	assert.Contains(t, view, "test-card")
	assert.Contains(t, view, "Card Number")
	assert.Contains(t, view, "16 or 18 digits")
	assert.Contains(t, view, "1234567890123456")
}

func TestUIController_addItemErrorView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			addItemErrorMsg: "Test error message",
		},
	}

	view := ui.addItemErrorView()

	assert.Contains(t, view, "Add Item Error")
	assert.Contains(t, view, "Test error message")
	assert.Contains(t, view, "Enter to return")
}

func TestUIController_addItemSuccessView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			addItemSuccessMsg: "Test success message",
		},
	}

	view := ui.addItemSuccessView()

	assert.Contains(t, view, "Item Added Successfully")
	assert.Contains(t, view, "Test success message")
	assert.Contains(t, view, "Enter to return")
}
