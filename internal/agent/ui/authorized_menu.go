package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) menuLoggedInView() string {
	title := titleStyle.Render(fmt.Sprintf("Welcome, %s!", ui.login))
	subtitle := "Choose an option - enter number or use arrow keys:"

	options := []string{
		"View Items",
		"Add Item",
		"Logout",
	}

	menu := ""
	for i, option := range options {
		prefix := fmt.Sprintf("%d. ", i+1)
		if i == ui.loggedInMenu {
			menu += selectedStyle.Render(prefix+option) + "\n"
		} else {
			menu += menuStyle.Render(prefix+option) + "\n"
		}
	}

	controls := "\nControls: ↑/↓ to navigate, Enter to select, q to quit"
	return fmt.Sprintf("%s\n\n%s\n\n%s%s", title, subtitle, menu, controls)
}

func (ui *UIController) handleMenuLoggedInInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "up", "k":
		if ui.loggedInMenu > 0 {
			ui.loggedInMenu--
		}
	case "down", "j":
		if ui.loggedInMenu < 2 {
			ui.loggedInMenu++
		}
	case "1":
		ui.loggedInMenu = 0
		return ui.handleViewItems()
	case "2":
		ui.loggedInMenu = 1
		return ui.handleAddItem()
	case "3":
		ui.loggedInMenu = 2
		return ui.handleLogout()
	case "enter":
		switch ui.loggedInMenu {
		case 0:
			return ui.handleViewItems()
		case 1:
			return ui.handleAddItem()
		case 2:
			return ui.handleLogout()
		}
	}
	return ui, nil
}