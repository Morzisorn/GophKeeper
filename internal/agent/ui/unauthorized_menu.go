package ui

import "fmt"

func (ui *UIController) menuLoggedOutView() string {
	title := titleStyle.Render("Welcome to GophKeeper - best sensitive info manager :)")
	subtitle := "Choose an option - enter number or use arrow keys:"

	// Show authentication error if present
	errorMsg := ""
	if authError := ui.messages.Get("auth_error"); authError != "" {
		errorMsg = "\n" + errorStyle.Render("❌ "+authError) + "\n"
	}

	options := []string{
		"Sign up",
		"Sign in",
	}

	menu := ""
	for i, option := range options {
		prefix := fmt.Sprintf("%d. ", i+1)
		if i == ui.currentMenu {
			menu += selectedStyle.Render(prefix+option) + "\n"
		} else {
			menu += menuStyle.Render(prefix+option) + "\n"
		}
	}

	controls := "\nControls: ↑/↓ to navigate, Enter to select, q to quit"

	return fmt.Sprintf("%s\n\n%s%s\n\n%s%s", title, subtitle, errorMsg, menu, controls)
}
