package ui

import (
	"gophkeeper/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_deleteItemCmd(t *testing.T) {
	// Этот тест сложен для unit-тестирования, так как требует mock Item service
	// Пропускаем, так как это интеграционный тест
	t.Skip("Requires mock Item service - integration test")
}

func TestUIController_handleConfirmDeleteInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleConfirmDeleteInput_Escape(t *testing.T) {
	ui := &UIController{
		state: stateConfirmDelete,
	}

	model, cmd := ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state)
}

func TestUIController_handleConfirmDeleteInput_Navigation_Left(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
	}

	model, cmd := ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyLeft})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.confirmChoice)

	// Test h key
	ui.confirmChoice = 1
	model, cmd = ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.confirmChoice)
}

func TestUIController_handleConfirmDeleteInput_Navigation_Right(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
	}

	model, cmd := ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyRight})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.confirmChoice)

	// Test l key
	ui.confirmChoice = 0
	model, cmd = ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.confirmChoice)
}

func TestUIController_handleConfirmDeleteInput_QuickChoice_Yes(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{Name: "test-item"},
			},
		},
	}

	model, cmd := ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return delete command
	assert.Equal(t, 1, ui.confirmChoice)
	assert.Equal(t, stateProcessing, ui.state)
}

func TestUIController_handleConfirmDeleteInput_QuickChoice_No(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
		state:         stateConfirmDelete,
	}

	model, cmd := ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.confirmChoice)
	assert.Equal(t, stateItemDetails, ui.state)
}

func TestUIController_handleConfirmDeleteInput_Enter_Confirm(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1, // Yes selected
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{Name: "test-item"},
			},
		},
	}

	model, cmd := ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return delete command
	assert.Equal(t, stateProcessing, ui.state)
}

func TestUIController_handleConfirmDeleteInput_Enter_Cancel(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0, // No selected
		state:         stateConfirmDelete,
	}

	model, cmd := ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state)
}

func TestUIController_handleConfirmDelete_Confirm_ValidItem(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
		itemCtrl: itemCtrl{
			currentItem: 0,
			items: []models.EncryptedItem{
				{ID: [16]byte{1, 2, 3}, Name: "test-item"},
			},
		},
	}

	result, cmd := ui.handleConfirmDelete()

	assert.Equal(t, ui, result)
	assert.NotNil(t, cmd) // Should return delete command
	assert.Equal(t, stateProcessing, ui.state)
}

func TestUIController_handleConfirmDelete_Cancel(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
		state:         stateConfirmDelete,
	}

	result, cmd := ui.handleConfirmDelete()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state)
}

func TestUIController_handleConfirmDelete_InvalidIndex(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
		itemCtrl: itemCtrl{
			currentItem: 5, // Index out of bounds
			items: []models.EncryptedItem{
				{Name: "test-item"},
			},
		},
		state: stateConfirmDelete,
	}

	result, cmd := ui.handleConfirmDelete()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemDetails, ui.state)
}

func TestUIController_handleDeleteErrorInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleDeleteErrorInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleDeleteErrorInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleDeleteErrorInput_Enter(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			deleteErrorMsg: "test error message",
		},
	}

	model, cmd := ui.handleDeleteErrorInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateConfirmDelete, ui.state)
	assert.Empty(t, ui.itemCtrl.deleteErrorMsg)
}

func TestUIController_handleDeleteErrorInput_Escape(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			deleteErrorMsg: "test error message",
		},
	}

	model, cmd := ui.handleDeleteErrorInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateConfirmDelete, ui.state)
	assert.Empty(t, ui.itemCtrl.deleteErrorMsg)
}

func TestUIController_handleDeleteSuccessInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleDeleteSuccessInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleDeleteSuccessInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleDeleteSuccessInput_Enter(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			deleteSuccessMsg: "test success message",
		},
	}

	model, cmd := ui.handleDeleteSuccessInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return loadItemsCmd
	assert.Equal(t, stateProcessing, ui.state)
	assert.Empty(t, ui.itemCtrl.deleteSuccessMsg)
}

func TestUIController_deleteSuccessView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			deleteSuccessMsg: "Item deleted successfully!",
		},
	}

	view := ui.deleteSuccessView()

	assert.Contains(t, view, "Item Deleted Successfully")
	assert.Contains(t, view, "Item deleted successfully!")
	assert.Contains(t, view, "Enter to return to items list")
	assert.Contains(t, view, "q to quit")
}

func TestUIController_deleteErrorView(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			deleteErrorMsg: "Failed to delete item: server error",
		},
	}

	view := ui.deleteErrorView()

	assert.Contains(t, view, "Delete Error")
	assert.Contains(t, view, "Failed to delete item: server error")
	assert.Contains(t, view, "Enter to try again")
	assert.Contains(t, view, "Esc to cancel")
	assert.Contains(t, view, "q to quit")
}

func TestUIController_confirmDeleteView_ValidItem(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0, // No selected
		itemCtrl: itemCtrl{
			selectedItem: &models.EncryptedItem{
				Name: "Test Item",
			},
		},
	}

	view := ui.confirmDeleteView()

	assert.Contains(t, view, "Confirm Delete")
	assert.Contains(t, view, "Are you sure you want to delete 'Test Item'?")
	assert.Contains(t, view, "[ No ]")
	assert.Contains(t, view, "[ Yes ]")
	assert.Contains(t, view, "←/→ or h/l to navigate")
	assert.Contains(t, view, "y/n for quick choice")
}

func TestUIController_confirmDeleteView_YesSelected(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1, // Yes selected
		itemCtrl: itemCtrl{
			selectedItem: &models.EncryptedItem{
				Name: "Test Item",
			},
		},
	}

	view := ui.confirmDeleteView()

	assert.Contains(t, view, "Confirm Delete")
	assert.Contains(t, view, "Test Item")
	// В этом случае "Yes" должен быть выделен, а "No" - обычным
	assert.Contains(t, view, "[ No ]")
	assert.Contains(t, view, "[ Yes ]")
}

func TestUIController_confirmDeleteView_NoItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			selectedItem: nil,
		},
	}

	view := ui.confirmDeleteView()

	assert.Equal(t, "No item selected", view)
}

func TestUIController_handleConfirmDeleteInput_OtherKey(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
		state:         stateConfirmDelete,
	}

	model, cmd := ui.handleConfirmDeleteInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.confirmChoice)          // Should remain unchanged
	assert.Equal(t, stateConfirmDelete, ui.state) // Should remain unchanged
}
