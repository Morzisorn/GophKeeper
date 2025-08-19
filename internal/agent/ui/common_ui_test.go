package ui

import (
	"gophkeeper/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_Update_KeyMsg(t *testing.T) {
	ui := &UIController{
		state: stateMenuLoggedOut,
	}

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	model, cmd := ui.Update(keyMsg)

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return quit command
}

func TestUIController_Update_ItemDecrypted(t *testing.T) {
	ui := &UIController{}
	testItem := &models.Item{Name: "test-item"}

	msg := itemDecrypted{item: testItem}
	model, cmd := ui.Update(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, testItem, ui.itemCtrl.decryptedItem)
}

func TestUIController_Update_ProcessComplete(t *testing.T) {
	ui := &UIController{}

	msg := processComplete{
		success: true,
		message: "test message",
		context: "auth",
	}

	model, cmd := ui.Update(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
}

func TestUIController_Update_ItemsLoaded(t *testing.T) {
	ui := &UIController{}
	testItems := []models.EncryptedItem{
		{Name: "item1"},
		{Name: "item2"},
	}

	msg := itemsLoaded{items: testItems}
	model, cmd := ui.Update(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, testItems, ui.itemCtrl.items)
	assert.Equal(t, 0, ui.itemCtrl.currentItem)
	assert.Equal(t, stateItemsList, ui.state)
}

func TestUIController_Update_ErrorMsg(t *testing.T) {
	ui := &UIController{
		state: stateProcessing,
	}

	msg := errorMsg{
		err:     assert.AnError,
		context: "test-context",
	}

	model, cmd := ui.Update(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateError, ui.state)
}

func TestUIController_View_MenuLoggedOut(t *testing.T) {
	ui := &UIController{
		state: stateMenuLoggedOut,
	}

	view := ui.View()

	assert.NotEmpty(t, view)
	assert.NotContains(t, view, "View error:")
}

func TestUIController_View_MenuLoggedIn(t *testing.T) {
	ui := &UIController{
		state: stateMenuLoggedIn,
		userCtrl: userCtrl{
			login: "test-user",
		},
	}

	view := ui.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Welcome, test-user!")
	assert.NotContains(t, view, "View error:")
}

func TestUIController_View_Processing(t *testing.T) {
	ui := &UIController{
		state: stateProcessing,
	}

	view := ui.View()

	assert.NotEmpty(t, view)
	assert.Contains(t, view, "Processing")
	assert.NotContains(t, view, "View error:")
}

func TestUIController_View_InvalidState(t *testing.T) {
	ui := &UIController{
		state: state(999), // Invalid state
	}

	view := ui.View()

	assert.Contains(t, view, "View error:")
	assert.Contains(t, view, "Current state: 999")
}

func TestUIController_successView(t *testing.T) {
	ui := &UIController{}
	ui.messages.Set("success", "Operation completed successfully")
	ui.messages.Set("success_context", "test-context")

	view := ui.successView()

	assert.Contains(t, view, "Success!")
	assert.Contains(t, view, "Operation completed successfully")
	assert.Contains(t, view, "[test-context]")
	assert.Contains(t, view, "Press Enter to continue")
}

func TestUIController_successView_NoContext(t *testing.T) {
	ui := &UIController{}
	ui.messages.Set("success", "Operation completed")

	view := ui.successView()

	assert.Contains(t, view, "Success!")
	assert.Contains(t, view, "Operation completed")
	assert.NotContains(t, view, "[")
	assert.Contains(t, view, "Press Enter to continue")
}

func TestUIController_errorView(t *testing.T) {
	ui := &UIController{}
	ui.messages.Set("error", "Something went wrong")
	ui.messages.Set("error_context", "error-context")

	view := ui.errorView()

	assert.Contains(t, view, "Error!")
	assert.Contains(t, view, "Something went wrong")
	assert.Contains(t, view, "[error-context]")
	assert.Contains(t, view, "Press Enter to continue, q to quit")
}

func TestUIController_errorView_NoContext(t *testing.T) {
	ui := &UIController{}
	ui.messages.Set("error", "Something failed")

	view := ui.errorView()

	assert.Contains(t, view, "Error!")
	assert.Contains(t, view, "Something failed")
	assert.NotContains(t, view, "[")
	assert.Contains(t, view, "Press Enter to continue, q to quit")
}

func TestUIController_handleResultInput_Quit(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleResultInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	model, cmd = ui.handleResultInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleResultInput_Enter_Success_Authenticated(t *testing.T) {
	ui := &UIController{
		state: stateSuccess,
		userCtrl: userCtrl{
			isAuthenticated: true,
		},
		loggedInMenu: 5,
	}

	model, cmd := ui.handleResultInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Equal(t, 0, ui.loggedInMenu)
}

func TestUIController_handleResultInput_Enter_Success_NotAuthenticated(t *testing.T) {
	ui := &UIController{
		state: stateSuccess,
		userCtrl: userCtrl{
			isAuthenticated: false,
		},
	}

	model, cmd := ui.handleResultInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should quit
	assert.Equal(t, stateFinished, ui.state)
}

func TestUIController_handleResultInput_Enter_Error_Authenticated(t *testing.T) {
	ui := &UIController{
		state: stateError,
		userCtrl: userCtrl{
			isAuthenticated: true,
		},
		loggedInMenu: 5,
	}

	model, cmd := ui.handleResultInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Equal(t, 0, ui.loggedInMenu)
}

func TestUIController_handleResultInput_Enter_Error_NotAuthenticated(t *testing.T) {
	ui := &UIController{
		state: stateError,
		userCtrl: userCtrl{
			isAuthenticated: false,
			login:           "test-user",
		},
		menuCtrl: menuCtrl{
			currentMenu: 5,
		},
	}

	model, cmd := ui.handleResultInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedOut, ui.state)
	assert.Equal(t, 0, ui.menuCtrl.currentMenu)
	assert.Empty(t, ui.userCtrl.login)
}

func TestUIController_handleProcessComplete_Success_AuthToMaster(t *testing.T) {
	ui := &UIController{
		input: "test-input",
	}

	msg := processComplete{
		success: true,
		context: "auth_to_master",
		message: "Auth successful",
	}

	model, cmd := ui.handleProcessComplete(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMasterPassword, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_handleProcessComplete_Success_MasterPassword(t *testing.T) {
	ui := &UIController{
		input: "test-input",
	}

	msg := processComplete{
		success: true,
		context: "master_password",
		message: "Master password set",
	}

	model, cmd := ui.handleProcessComplete(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_handleProcessComplete_Success_DeleteItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			selectedItem:  &models.EncryptedItem{Name: "test"},
			decryptedItem: &models.Item{Name: "test"},
		},
	}

	msg := processComplete{
		success: true,
		context: "delete_item",
		message: "Item deleted",
	}

	model, cmd := ui.handleProcessComplete(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateDeleteSuccess, ui.state)
	assert.Equal(t, "Item deleted", ui.itemCtrl.deleteSuccessMsg)
	assert.Nil(t, ui.itemCtrl.selectedItem)
	assert.Nil(t, ui.itemCtrl.decryptedItem)
}

func TestUIController_handleProcessComplete_Success_SaveItem(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			newItem: models.Item{Name: "test-item"},
		},
	}

	msg := processComplete{
		success: true,
		context: "save_item",
		message: "Item saved",
	}

	model, cmd := ui.handleProcessComplete(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateAddItemSuccess, ui.state)
	assert.Equal(t, "Item saved", ui.itemCtrl.addItemSuccessMsg)
	assert.Equal(t, models.Item{}, ui.itemCtrl.newItem) // Should be reset
}

func TestUIController_handleProcessComplete_Failure_Auth(t *testing.T) {
	ui := &UIController{
		input: "test-input",
	}

	msg := processComplete{
		success: false,
		context: "auth",
		message: "Authentication failed",
	}

	model, cmd := ui.handleProcessComplete(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedOut, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_handleProcessComplete_Failure_DeleteItem(t *testing.T) {
	ui := &UIController{}

	msg := processComplete{
		success: false,
		context: "delete_item",
		message: "Delete failed",
	}

	model, cmd := ui.handleProcessComplete(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateDeleteError, ui.state)
	assert.Equal(t, "Delete failed", ui.itemCtrl.deleteErrorMsg)
}

func TestUIController_handleProcessComplete_Failure_SaveItem(t *testing.T) {
	ui := &UIController{}

	msg := processComplete{
		success: false,
		context: "save_item",
		message: "Save failed",
	}

	model, cmd := ui.handleProcessComplete(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateAddItemError, ui.state)
	assert.Equal(t, "Save failed", ui.itemCtrl.addItemErrorMsg)
}
