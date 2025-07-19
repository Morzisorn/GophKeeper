package ui

import (
	"gophkeeper/models"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_handleLogout(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1, // Set to non-zero initially
	}

	result, cmd := ui.handleLogout()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, stateConfirmLogout, ui.state)
	assert.Equal(t, 0, ui.confirmChoice) // Should be reset to 0
}

func TestUIController_handleConfirmLogoutInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleConfirmLogoutInput_Escape(t *testing.T) {
	ui := &UIController{
		state:        stateConfirmLogout,
		loggedInMenu: 0,
	}

	model, cmd := ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Equal(t, 2, ui.loggedInMenu) // Should be set to logout option
}

func TestUIController_handleConfirmLogoutInput_Navigation_Left(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
	}

	model, cmd := ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyLeft})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.confirmChoice)

	// Test h key
	ui.confirmChoice = 1
	model, cmd = ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.confirmChoice)
}

func TestUIController_handleConfirmLogoutInput_Navigation_Right(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
	}

	model, cmd := ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyRight})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.confirmChoice)

	// Test l key
	ui.confirmChoice = 0
	model, cmd = ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.confirmChoice)
}

func TestUIController_handleConfirmLogoutInput_QuickChoice_Yes(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
	}

	model, cmd := ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return logout command
	assert.Equal(t, 1, ui.confirmChoice)
	assert.Equal(t, stateProcessing, ui.state)
}

func TestUIController_handleConfirmLogoutInput_QuickChoice_No(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
		state:         stateConfirmLogout,
		loggedInMenu:  0,
	}

	model, cmd := ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.confirmChoice)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Equal(t, 2, ui.loggedInMenu)
}

func TestUIController_handleConfirmLogoutInput_Enter_Confirm(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1, // Yes selected
	}

	model, cmd := ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // Should return logout command
	assert.Equal(t, stateProcessing, ui.state)
}

func TestUIController_handleConfirmLogoutInput_Enter_Cancel(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0, // No selected
		state:         stateConfirmLogout,
		loggedInMenu:  0,
	}

	model, cmd := ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Equal(t, 2, ui.loggedInMenu)
}

func TestUIController_executeLogout_Confirm(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1,
	}

	result, cmd := ui.executeLogout()

	assert.Equal(t, ui, result)
	assert.NotNil(t, cmd) // Should return logoutCmd
	assert.Equal(t, stateProcessing, ui.state)
}

func TestUIController_executeLogout_Cancel(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
		state:         stateConfirmLogout,
		loggedInMenu:  0,
	}

	result, cmd := ui.executeLogout()

	assert.Equal(t, ui, result)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Equal(t, 2, ui.loggedInMenu)
}

func TestUIController_logoutCmd(t *testing.T) {
	// Этот тест сложен для unit-тестирования, так как требует mock User service
	// Пропускаем, так как это интеграционный тест
	t.Skip("Requires mock User service - integration test")
}

func TestUIController_handleLogoutSuccessInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleLogoutSuccessInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleLogoutSuccessInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleLogoutSuccessInput_Enter(t *testing.T) {
	ui := &UIController{
		userCtrl: userCtrl{
			isAuthenticated: true,
			login:           "test-user",
		},
		menuCtrl: menuCtrl{
			currentMenu: 5,
		},
		logoutCtrl: logoutCtrl{
			logoutSuccessMsg: "test success message",
		},
		itemCtrl: itemCtrl{
			items: []models.EncryptedItem{{Name: "test"}},
		},
	}

	model, cmd := ui.handleLogoutSuccessInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedOut, ui.state)
	assert.Equal(t, 0, ui.menuCtrl.currentMenu)
	assert.Empty(t, ui.logoutCtrl.logoutSuccessMsg)

	// Check that user session is cleared
	assert.False(t, ui.userCtrl.isAuthenticated)
	assert.Empty(t, ui.userCtrl.login)
	assert.Nil(t, ui.itemCtrl.items)
}

func TestUIController_handleLogoutErrorInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleLogoutErrorInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleLogoutErrorInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleLogoutErrorInput_Enter(t *testing.T) {
	ui := &UIController{
		loggedInMenu: 0,
		logoutCtrl: logoutCtrl{
			logoutErrorMsg: "test error message",
		},
	}

	model, cmd := ui.handleLogoutErrorInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Equal(t, 2, ui.loggedInMenu)
	assert.Empty(t, ui.logoutCtrl.logoutErrorMsg)
}

func TestUIController_handleLogoutErrorInput_Escape(t *testing.T) {
	ui := &UIController{
		loggedInMenu: 0,
		logoutCtrl: logoutCtrl{
			logoutErrorMsg: "test error message",
		},
	}

	model, cmd := ui.handleLogoutErrorInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedIn, ui.state)
	assert.Equal(t, 2, ui.loggedInMenu)
	assert.Empty(t, ui.logoutCtrl.logoutErrorMsg)
}

func TestUIController_clearUserSession(t *testing.T) {
	ui := &UIController{
		userCtrl: userCtrl{
			isAuthenticated: true,
			login:           "test-user",
		},
		itemCtrl: itemCtrl{
			items:         []models.EncryptedItem{{Name: "test"}},
			currentItem:   5,
			selectedItem:  &models.EncryptedItem{Name: "selected"},
			decryptedItem: &models.Item{Name: "decrypted"},
			editingItem:   &models.Item{Name: "editing"},
			newItem:       models.Item{Name: "new"},
			itemMetaCtrl: itemMetaCtrl{
				metadataKeys:       []string{"key1", "key2"},
				currentMetaKey:     "current-key",
				currentMetaValue:   "current-value",
				metadataSuccessMsg: "meta success",
				metadataErrorMsg:   "meta error",
			},
			addItemErrorMsg:   "add error",
			addItemSuccessMsg: "add success",
			deleteSuccessMsg:  "delete success",
			deleteErrorMsg:    "delete error",
			editSuccessMsg:    "edit success",
			editErrorMsg:      "edit error",
			decryptErrorMsg:   "decrypt error",
		},
	}

	// Set some messages
	ui.messages.Set("info", "test info")
	ui.messages.Set("error", "test error")

	ui.clearUserSession()

	// Check user session is cleared
	assert.False(t, ui.userCtrl.isAuthenticated)
	assert.Empty(t, ui.userCtrl.login)

	// Check items are cleared
	assert.Nil(t, ui.itemCtrl.items)
	assert.Equal(t, 0, ui.itemCtrl.currentItem)
	assert.Nil(t, ui.itemCtrl.selectedItem)
	assert.Nil(t, ui.itemCtrl.decryptedItem)
	assert.Nil(t, ui.itemCtrl.editingItem)
	assert.Equal(t, models.Item{}, ui.itemCtrl.newItem)

	// Check metadata is cleared
	assert.Nil(t, ui.itemCtrl.itemMetaCtrl.metadataKeys)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.currentMetaKey)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.currentMetaValue)

	// Check all messages are cleared
	assert.Empty(t, ui.messages.Get("info"))
	assert.Empty(t, ui.messages.Get("error"))
	assert.Empty(t, ui.itemCtrl.addItemErrorMsg)
	assert.Empty(t, ui.itemCtrl.addItemSuccessMsg)
	assert.Empty(t, ui.itemCtrl.deleteSuccessMsg)
	assert.Empty(t, ui.itemCtrl.deleteErrorMsg)
	assert.Empty(t, ui.itemCtrl.editSuccessMsg)
	assert.Empty(t, ui.itemCtrl.editErrorMsg)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.metadataSuccessMsg)
	assert.Empty(t, ui.itemCtrl.itemMetaCtrl.metadataErrorMsg)
	assert.Empty(t, ui.itemCtrl.decryptErrorMsg)
}

func TestUIController_handleConfirmLogoutInput_OtherKey(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0,
		state:         stateConfirmLogout,
	}

	model, cmd := ui.handleConfirmLogoutInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.confirmChoice)          // Should remain unchanged
	assert.Equal(t, stateConfirmLogout, ui.state) // Should remain unchanged
}

// View tests
func TestUIController_confirmLogoutView_NoSelected(t *testing.T) {
	ui := &UIController{
		confirmChoice: 0, // No selected
	}

	view := ui.confirmLogoutView()

	assert.Contains(t, view, "Confirm Logout")
	assert.Contains(t, view, "Are you sure you want to logout?")
	assert.Contains(t, view, "[ No ]")
	assert.Contains(t, view, "[ Yes ]")
	assert.Contains(t, view, "←/→ or h/l to navigate")
	assert.Contains(t, view, "y/n for quick choice")
	assert.Contains(t, view, "Enter to confirm")
	assert.Contains(t, view, "Esc to cancel")
}

func TestUIController_confirmLogoutView_YesSelected(t *testing.T) {
	ui := &UIController{
		confirmChoice: 1, // Yes selected
	}

	view := ui.confirmLogoutView()

	assert.Contains(t, view, "Confirm Logout")
	assert.Contains(t, view, "Are you sure you want to logout?")
	assert.Contains(t, view, "[ No ]")
	assert.Contains(t, view, "[ Yes ]")
	// В этом случае "Yes" должен быть выделен, а "No" - обычным
}

func TestUIController_logoutSuccessView(t *testing.T) {
	ui := &UIController{
		logoutCtrl: logoutCtrl{
			logoutSuccessMsg: "Successfully logged out!",
		},
	}

	view := ui.logoutSuccessView()

	assert.Contains(t, view, "Logout Successful")
	assert.Contains(t, view, "Successfully logged out!")
	assert.Contains(t, view, "Enter to continue")
	assert.Contains(t, view, "q to quit")
}

func TestUIController_logoutErrorView(t *testing.T) {
	ui := &UIController{
		logoutCtrl: logoutCtrl{
			logoutErrorMsg: "Failed to logout: server error",
		},
	}

	view := ui.logoutErrorView()

	assert.Contains(t, view, "Logout Error")
	assert.Contains(t, view, "Failed to logout: server error")
	assert.Contains(t, view, "Enter to return to menu")
	assert.Contains(t, view, "Esc to cancel")
	assert.Contains(t, view, "q to quit")
}
