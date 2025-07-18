package ui

import (
	"fmt"
	"gophkeeper/models"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) handleViewItemDetails() (*UIController, tea.Cmd) {
	if ui.currentItem < len(ui.items) {
		ui.selectedItem = &ui.items[ui.currentItem]
		ui.state = stateItemDetails
		ui.input = ""
		ui.decryptedItem = nil 

		return ui, ui.decryptItemCmd(ui.selectedItem)
	}
	return ui, nil
}

func (ui *UIController) handleItemDetailsInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc", "b":
		ui.state = stateItemsList
		return ui, nil
	case "d":
		ui.state = stateConfirmDelete
		ui.confirmChoice = 0
		return ui, nil
	case "e":
		if ui.decryptedItem != nil {
			return ui.startEditItem()
		}
		return ui, nil
	case "m":
		if ui.decryptedItem != nil {
			return ui.startManageMetadata()
		}
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) itemDetailsView() string {
	if ui.currentItem >= len(ui.items) {
		return "No item selected"
	}

	selectedItem := &ui.items[ui.currentItem]
	title := titleStyle.Render(fmt.Sprintf("Item Details: %s", selectedItem.Name))

	details := fmt.Sprintf("Type: %s\n", selectedItem.Type)
	details += fmt.Sprintf("Created: %s\n", selectedItem.CreatedAt.Format("2006-01-02 15:04:05"))
	details += fmt.Sprintf("Updated: %s\n\n", selectedItem.UpdatedAt.Format("2006-01-02 15:04:05"))

	if ui.decryptedItem != nil {
		details += "Data:\n"
		switch data := ui.decryptedItem.Data.(type) {
		case *models.Credentials:
			details += fmt.Sprintf("  Login: %s\n", data.Login)
			details += fmt.Sprintf("  Password: %s\n", data.Password)
		case *models.Text:
			details += fmt.Sprintf("  Content: %s\n", data.Content)
		case *models.Card:
			details += fmt.Sprintf("  Number: %s\n", data.Number)
			details += fmt.Sprintf("  Expiry: %s\n", data.ExpiryDate)
			details += fmt.Sprintf("  CVV: %s\n", data.SecurityCode)
			details += fmt.Sprintf("  Cardholder: %s\n", data.CardholderName)
		case *models.Binary:
			details += fmt.Sprintf("  Content: %s\n", string(data.Content))
		default:
			details += "  Unknown data type\n"
		}

		if len(ui.decryptedItem.Meta.Map) > 0 {
			details += "\nMetadata:\n"
			for key, value := range ui.decryptedItem.Meta.Map {
				details += fmt.Sprintf("  %s: %s\n", key, value)
			}
		} else {
			details += "\nNo metadata\n"
		}
	} else {
		details += "Loading data...\n"
	}

	controls := "\nControls: e to edit, m to manage metadata, d to delete, b/Esc to go back"
	return fmt.Sprintf("%s\n\n%s%s", title, details, controls)
}

