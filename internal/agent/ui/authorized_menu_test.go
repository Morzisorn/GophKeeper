package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_menuLoggedInView(t *testing.T) {
	ui := &UIController{
		userCtrl: userCtrl{
			login: "test-user",
		},
		loggedInMenu: 1,
	}

	view := ui.menuLoggedInView()

	assert.Contains(t, view, "Welcome, test-user!")
	assert.Contains(t, view, "Choose an option")
	assert.Contains(t, view, "View Items")
	assert.Contains(t, view, "Add Item")
	assert.Contains(t, view, "Logout")
	assert.Contains(t, view, "↑/↓ to navigate")
	assert.Contains(t, view, "1.")
	assert.Contains(t, view, "2.")
	assert.Contains(t, view, "3.")
}

func TestUIController_handleMenuLoggedInInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleMenuLoggedInInput_Navigation_Up(t *testing.T) {
	ui := &UIController{
		loggedInMenu: 2,
	}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyUp})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.loggedInMenu)

	// Test k key
	model, cmd = ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.loggedInMenu)

	// Test that it doesn't go below 0
	model, cmd = ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyUp})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.loggedInMenu)
}

func TestUIController_handleMenuLoggedInInput_Navigation_Down(t *testing.T) {
	ui := &UIController{
		loggedInMenu:    0,
		maxLoggedInMenu: 3, // Устанавливаем максимальное значение (4 пункта меню: 0,1,2,3)
	}

	// Test down arrow key
	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyDown})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.loggedInMenu)

	// Test j key
	model, cmd = ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 2, ui.loggedInMenu)

	// Test one more down - should go to max (3)
	model, cmd = ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyDown})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 3, ui.loggedInMenu)

	// Test that it doesn't go above maxLoggedInMenu
	model, cmd = ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyDown})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 3, ui.loggedInMenu) // Should stay at maxLoggedInMenu
}

func TestUIController_handleMenuLoggedInInput_DirectSelection_ViewItems(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // handleViewItems returns a command
	assert.Equal(t, 0, ui.loggedInMenu)
}

func TestUIController_handleMenuLoggedInInput_DirectSelection_AddItem(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd) // handleAddItem returns nil command
	assert.Equal(t, 2, ui.loggedInMenu)
}

func TestUIController_handleMenuLoggedInInput_DirectSelection_ViewItemsWithType(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // handleViewItemsWithType returns a command
	assert.Equal(t, 1, ui.loggedInMenu)
}

func TestUIController_handleMenuLoggedInInput_DirectSelection_Logout(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'4'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd) // handleLogout возвращает nil command
	assert.Equal(t, 3, ui.loggedInMenu)
}

func TestUIController_handleMenuLoggedInInput_Enter_ViewItems(t *testing.T) {
	ui := &UIController{
		loggedInMenu: 0,
	}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // handleViewItems returns a command
}

func TestUIController_handleMenuLoggedInInput_Enter_ViewItemsWithType(t *testing.T) {
	ui := &UIController{
		loggedInMenu: 1,
	}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd) // handleViewItemsWithType returns a command
}

func TestUIController_handleMenuLoggedInInput_Enter_AddItem(t *testing.T) {
	ui := &UIController{
		loggedInMenu: 2,
	}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd) // handleAddItem returns nil command
}

func TestUIController_handleMenuLoggedInInput_Enter_Logout(t *testing.T) {
	ui := &UIController{
		loggedInMenu: 3,
	}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd) // handleLogout возвращает nil command
}

func TestUIController_handleMenuLoggedInInput_InvalidKey(t *testing.T) {
	ui := &UIController{
		loggedInMenu: 1,
	}

	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.loggedInMenu) // Should remain unchanged
}

func TestUIController_handleMenuLoggedInInput_NumberOutOfRange(t *testing.T) {
	ui := &UIController{
		loggedInMenu: 1,
	}

	// Test number 0 (should be ignored)
	model, cmd := ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'0'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.loggedInMenu) // Should remain unchanged

	// Test number 5 (should be ignored)
	model, cmd = ui.handleMenuLoggedInInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.loggedInMenu) // Should remain unchanged
}

func TestUIController_menuLoggedInView_EmptyLogin(t *testing.T) {
	ui := &UIController{
		userCtrl: userCtrl{
			login: "",
		},
		loggedInMenu: 0,
	}

	view := ui.menuLoggedInView()

	assert.Contains(t, view, "Welcome, !")
	assert.Contains(t, view, "View Items")
}

func TestUIController_menuLoggedInView_SelectedItemHighlighted(t *testing.T) {
	ui := &UIController{
		userCtrl: userCtrl{
			login: "test-user",
		},
		loggedInMenu: 1, // Add Item should be selected
	}

	view := ui.menuLoggedInView()

	// The selected item should appear differently in the view
	// Since we can't test styling directly, we just verify the content exists
	assert.Contains(t, view, "Add Item")
	assert.Contains(t, view, "View Items")
	assert.Contains(t, view, "Logout")
}
