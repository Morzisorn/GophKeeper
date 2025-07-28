package ui

import (
	"fmt"
	"gophkeeper/models"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) handleLogout() (tea.Model, tea.Cmd) {
	ui.state = stateConfirmLogout
	ui.confirmChoice = 0
	return ui, nil
}

func (ui *UIController) handleConfirmLogoutInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateMenuLoggedIn
		ui.loggedInMenu = 2
		return ui, nil
	case "left", "h":
		ui.confirmChoice = 0
	case "right", "l":
		ui.confirmChoice = 1
	case "y":
		ui.confirmChoice = 1
		return ui.executeLogout()
	case "n":
		ui.confirmChoice = 0
		ui.state = stateMenuLoggedIn
		ui.loggedInMenu = 2
		return ui, nil
	case "enter":
		return ui.executeLogout()
	}
	return ui, nil
}

func (ui *UIController) executeLogout() (*UIController, tea.Cmd) {
	if ui.confirmChoice == 1 {
		ui.state = stateProcessing
		return ui, ui.logoutCmd()
	} else {
		ui.state = stateMenuLoggedIn
		ui.loggedInMenu = 2
		return ui, nil
	}
}

func (ui *UIController) logoutCmd() tea.Cmd {
	return func() tea.Msg {
		err := ui.User.Logout()
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Logout error: %v", err),
				context: "logout",
			}
		}

		return processComplete{
			success: true,
			message: "Successfully logged out",
			context: "logout",
		}
	}
}

func (ui *UIController) confirmLogoutView() string {
	title := titleStyle.Render("Confirm Logout")
	warning := "Are you sure you want to logout?"

	options := ""
	if ui.confirmChoice == 0 {
		options += selectedStyle.Render("[ No ]") + "  "
		options += menuStyle.Render("[ Yes ]")
	} else {
		options += menuStyle.Render("[ No ]") + "  "
		options += selectedStyle.Render("[ Yes ]")
	}

	controls := "\nControls: ←/→ or h/l to navigate, y/n for quick choice, Enter to confirm, Esc to cancel"
	return fmt.Sprintf("%s\n\n%s\n\n%s%s", title, warning, options, controls)
}

func (ui *UIController) handleLogoutSuccessInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.clearUserSession()
		ui.state = stateMenuLoggedOut
		ui.currentMenu = 0
		ui.logoutSuccessMsg = ""
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) handleLogoutErrorInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.state = stateMenuLoggedIn
		ui.loggedInMenu = 2
		ui.logoutErrorMsg = ""
		return ui, nil
	case "esc":
		ui.state = stateMenuLoggedIn
		ui.loggedInMenu = 2
		ui.logoutErrorMsg = ""
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) logoutSuccessView() string {
	title := successStyle.Render("Logout Successful")
	message := ui.logoutSuccessMsg

	controls := "\nControls: Enter to continue, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}

func (ui *UIController) logoutErrorView() string {
	title := errorStyle.Render("Logout Error")
	message := ui.logoutErrorMsg

	controls := "\nControls: Enter to return to menu, Esc to cancel, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}

func (ui *UIController) clearUserSession() {
	ui.isAuthenticated = false
	ui.login = ""
	ui.items = nil
	ui.currentItem = 0
	ui.selectedItem = nil
	ui.decryptedItem = nil
	ui.editingItem = nil
	ui.newItem = models.Item{}
	ui.metadataKeys = nil
	ui.currentMetaKey = ""
	ui.currentMetaValue = ""
	ui.messages.ClearAll()

	ui.addItemErrorMsg = ""
	ui.addItemSuccessMsg = ""
	ui.deleteSuccessMsg = ""
	ui.deleteErrorMsg = ""
	ui.editSuccessMsg = ""
	ui.editErrorMsg = ""
	ui.metadataSuccessMsg = ""
	ui.metadataErrorMsg = ""
	ui.decryptErrorMsg = ""
}
