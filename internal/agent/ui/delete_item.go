package ui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) deleteItemCmd(itemID [16]byte) tea.Cmd {
	return func() tea.Msg {
		err := ui.Item.DeleteItem(context.Background(), ui.login, itemID)
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Delete error: %v", err),
				context: "delete_item",
			}
		}

		return processComplete{
			success: true,
			message: "Item deleted successfully",
			context: "delete_item",
		}
	}
}

func (ui *UIController) handleConfirmDeleteInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateItemDetails
		return ui, nil
	case "left", "h":
		ui.confirmChoice = 0 
	case "right", "l":
		ui.confirmChoice = 1 
	case "y":
		ui.confirmChoice = 1
		return ui.handleConfirmDelete()
	case "n":
		ui.confirmChoice = 0
		ui.state = stateItemDetails
		return ui, nil
	case "enter":
		return ui.handleConfirmDelete()
	}
	return ui, nil
}

func (ui *UIController) handleConfirmDelete() (*UIController, tea.Cmd) {
	if ui.confirmChoice == 1 && ui.currentItem < len(ui.items) {
		selectedItem := &ui.items[ui.currentItem]
		ui.state = stateProcessing
		return ui, ui.deleteItemCmd(selectedItem.ID)
	} else {
		ui.state = stateItemDetails
		return ui, nil
	}
}

func (ui *UIController) handleDeleteErrorInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter", "esc":
		ui.state = stateConfirmDelete
		ui.deleteErrorMsg = ""
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) handleDeleteSuccessInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.state = stateProcessing
		ui.deleteSuccessMsg = ""
		return ui, ui.loadItemsCmd()
	}
	return ui, nil
}

func (ui *UIController) deleteSuccessView() string {
	title := successStyle.Render("Item Deleted Successfully")
	message := ui.deleteSuccessMsg

	controls := "\nControls: Enter to return to items list, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}

func (ui *UIController) deleteErrorView() string {
	title := errorStyle.Render("Delete Error")
	message := ui.deleteErrorMsg

	controls := "\nControls: Enter to try again, Esc to cancel, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}

func (ui *UIController) confirmDeleteView() string {
	if ui.selectedItem == nil {
		return "No item selected"
	}

	title := titleStyle.Render("Confirm Delete")
	warning := fmt.Sprintf("Are you sure you want to delete '%s'?", ui.selectedItem.Name)

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
