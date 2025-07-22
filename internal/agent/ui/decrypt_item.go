package ui

import (
	"fmt"
	"gophkeeper/models"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type itemDecrypted struct {
	item *models.Item
}

type decryptError struct {
	err     error
	context string
}

func (ui *UIController) decryptItemCmd(item *models.EncryptedItem) tea.Cmd {
	return func() tea.Msg {
		decryptedItem, err := ui.Item.DecryptItem(item)
		if err != nil {
			errStr := err.Error()
			if strings.Contains(errStr, "failed to decrypt") ||
				strings.Contains(errStr, "cipher: message authentication failed") ||
				strings.Contains(errStr, "authentication failed") {
				return decryptError{
					err:     fmt.Errorf("incorrect master password"),
					context: "decrypt_item",
				}
			}

			return decryptError{
				err:     err,
				context: "decrypt_item",
			}
		}

		return itemDecrypted{
			item: decryptedItem,
		}
	}
}

func (ui *UIController) handleDecryptError(msg decryptError) (tea.Model, tea.Cmd) {
	switch msg.context {
	case "decrypt_item":
		if strings.Contains(msg.err.Error(), "incorrect master password") {
			ui.state = stateDecryptError
			ui.decryptErrorMsg = "Incorrect master password. Please restart the application and enter the correct master password."
		} else {
			ui.state = stateDecryptError
			ui.decryptErrorMsg = fmt.Sprintf("Failed to decrypt item: %v", msg.err)
		}
		return ui, nil
	default:
		ui.state = stateItemsList
		return ui, nil
	}
}

func (ui *UIController) handleDecryptErrorInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.state = stateMasterPassword
		ui.input = ""
		ui.decryptErrorMsg = ""
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) decryptErrorView() string {
	title := errorStyle.Render("Decryption Error")
	message := ui.decryptErrorMsg

	if strings.Contains(message, "incorrect master password") {
		controls := "\nControls: Enter to re-enter master password, q to quit"
		return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
	} else {
		controls := "\nControls: Enter to re-enter master password, q to quit"
		return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
	}
}

