package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) startManageMetadata() (*UIController, tea.Cmd) {
	if ui.decryptedItem == nil {
		return ui, nil
	}

	ui.state = stateMetadataList
	ui.metadataMenu = 0
	ui.currentMetaKey = ""
	ui.currentMetaValue = ""
	ui.metadataKeys = make([]string, 0, len(ui.decryptedItem.Meta.Map))

	for key := range ui.decryptedItem.Meta.Map {
		ui.metadataKeys = append(ui.metadataKeys, key)
	}

	return ui, nil
}

func (ui *UIController) handleMetadataListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc", "b":
		ui.state = stateItemDetails
		return ui, nil
	case "up", "k":
		if ui.metadataMenu > 0 {
			ui.metadataMenu--
		}
	case "down", "j":
		maxMenu := len(ui.metadataKeys)
		if ui.metadataMenu < maxMenu {
			ui.metadataMenu++
		}
	case "enter":
		if ui.metadataMenu < len(ui.metadataKeys) {
			key := ui.metadataKeys[ui.metadataMenu]
			return ui.startEditMetadata(key)
		} else {
			return ui.startAddMetadata()
		}
	case "d":
		if ui.metadataMenu < len(ui.metadataKeys) {
			key := ui.metadataKeys[ui.metadataMenu]
			return ui.startDeleteMetadata(key)
		}
		return ui, nil
	case "a":
		return ui.startAddMetadata()
	}
	return ui, nil
}

func (ui *UIController) metadataListView() string {
	if ui.decryptedItem == nil {
		return "No item selected"
	}

	title := titleStyle.Render(fmt.Sprintf("Manage Metadata: %s", ui.decryptedItem.Name))

	content := ""

	if len(ui.metadataKeys) == 0 {
		content += "No metadata found.\n\n"
	} else {
		content += "Current metadata:\n\n"
		for i, key := range ui.metadataKeys {
			value := ui.decryptedItem.Meta.Map[key]
			line := fmt.Sprintf("%s: %s", key, value)
			if i == ui.metadataMenu {
				content += selectedStyle.Render("→ "+line) + "\n"
			} else {
				content += menuStyle.Render("  "+line) + "\n"
			}
		}
		content += "\n"
	}

	addOption := "+ Add new metadata"
	if ui.metadataMenu == len(ui.metadataKeys) {
		content += selectedStyle.Render("→ "+addOption) + "\n"
	} else {
		content += menuStyle.Render("  "+addOption) + "\n"
	}

	controls := "\nControls: ↑/↓ to navigate, Enter to edit/add, d to delete, a to add, b/Esc to go back"
	return fmt.Sprintf("%s\n\n%s%s", title, content, controls)
}

func (ui *UIController) startAddMetadata() (*UIController, tea.Cmd) {
	ui.state = stateAddMetadataKey
	ui.input = ""
	ui.currentMetaKey = ""
	ui.currentMetaValue = ""
	return ui, nil
}

func (ui *UIController) handleAddMetadataKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateMetadataList
		ui.input = ""
		return ui, nil
	case "enter":
		key := strings.TrimSpace(ui.input)
		if key == "" {
			return ui, nil
		}

		if _, exists := ui.decryptedItem.Meta.Map[key]; exists {
			ui.messages.Set("error", "Metadata key already exists")
			return ui, nil
		}

		ui.currentMetaKey = key
		ui.input = ""
		ui.state = stateAddMetadataValue
		ui.messages.Clear("error")
		return ui, nil
	case "backspace":
		if len(ui.input) > 0 {
			ui.input = ui.input[:len(ui.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			ui.input += msg.String()
		}
	}
	return ui, nil
}

func (ui *UIController) handleAddMetadataValueInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateAddMetadataKey
		ui.input = ui.currentMetaKey
		return ui, nil
	case "enter":
		value := strings.TrimSpace(ui.input)
		if value == "" {
			return ui, nil
		}

		ui.currentMetaValue = value
		return ui, ui.saveMetadataCmd("add", ui.currentMetaKey, ui.currentMetaValue)
	case "backspace":
		if len(ui.input) > 0 {
			ui.input = ui.input[:len(ui.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			ui.input += msg.String()
		}
	}
	return ui, nil
}

func (ui *UIController) addMetadataKeyView() string {
	title := titleStyle.Render("Add Metadata - Enter Key")
	input := inputStyle.Render(ui.input + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to continue"
	return fmt.Sprintf("%s\n\nKey: %s%s%s", title, input, errorMsg, controls)
}

func (ui *UIController) addMetadataValueView() string {
	title := titleStyle.Render("Add Metadata - Enter Value")
	input := inputStyle.Render(ui.input + "█")

	controls := "\nControls: Esc to go back, Enter to save"
	return fmt.Sprintf("%s\n\nKey: %s\nValue: %s%s", title, ui.currentMetaKey, input, controls)
}

func (ui *UIController) startEditMetadata(key string) (*UIController, tea.Cmd) {
	ui.currentMetaKey = key
	ui.currentMetaValue = ui.decryptedItem.Meta.Map[key]
	ui.input = ui.currentMetaValue
	ui.state = stateEditMetadataValue
	return ui, nil
}

func (ui *UIController) handleEditMetadataValueInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateMetadataList
		ui.input = ""
		return ui, nil
	case "enter":
		value := strings.TrimSpace(ui.input)
		if value == "" {
			return ui, nil
		}

		return ui, ui.saveMetadataCmd("edit", ui.currentMetaKey, value)
	case "backspace":
		if len(ui.input) > 0 {
			ui.input = ui.input[:len(ui.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			ui.input += msg.String()
		}
	}
	return ui, nil
}

func (ui *UIController) editMetadataValueView() string {
	title := titleStyle.Render("Edit Metadata - Edit Value")
	input := inputStyle.Render(ui.input + "█")

	controls := "\nControls: Esc to cancel, Enter to save"
	return fmt.Sprintf("%s\n\nKey: %s\nValue: %s%s", title, ui.currentMetaKey, input, controls)
}

func (ui *UIController) startDeleteMetadata(key string) (*UIController, tea.Cmd) {
	ui.currentMetaKey = key
	ui.state = stateConfirmDeleteMetadata
	ui.confirmChoice = 0 
	return ui, nil
}

func (ui *UIController) handleConfirmDeleteMetadataInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateMetadataList
		return ui, nil
	case "left", "h":
		ui.confirmChoice = 0 
	case "right", "l":
		ui.confirmChoice = 1 
	case "y":
		ui.confirmChoice = 1
		return ui.handleConfirmDeleteMetadata()
	case "n":
		ui.confirmChoice = 0
		ui.state = stateMetadataList
		return ui, nil
	case "enter":
		return ui.handleConfirmDeleteMetadata()
	}
	return ui, nil
}

func (ui *UIController) handleConfirmDeleteMetadata() (*UIController, tea.Cmd) {
	if ui.confirmChoice == 1 {
		return ui, ui.saveMetadataCmd("delete", ui.currentMetaKey, "")
	} else {
		ui.state = stateMetadataList
		return ui, nil
	}
}

func (ui *UIController) confirmDeleteMetadataView() string {
	title := titleStyle.Render("Confirm Delete Metadata")
	warning := fmt.Sprintf("Are you sure you want to delete metadata key '%s'?", ui.currentMetaKey)

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

func (ui *UIController) saveMetadataCmd(action, key, value string) tea.Cmd {
	return func() tea.Msg {
		if ui.decryptedItem.Meta.Map == nil {
			ui.decryptedItem.Meta.Map = make(map[string]string)
		}

		switch action {
		case "add", "edit":
			ui.decryptedItem.Meta.Map[key] = value
		case "delete":
			delete(ui.decryptedItem.Meta.Map, key)
		}

		ui.decryptedItem.UpdatedAt = time.Now()

		err := ui.Item.EditItem(context.Background(), ui.decryptedItem)
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Failed to save metadata: %v", err),
				context: "save_metadata",
			}
		}

		var message string
		switch action {
		case "add":
			message = fmt.Sprintf("Metadata '%s' added successfully", key)
		case "edit":
			message = fmt.Sprintf("Metadata '%s' updated successfully", key)
		case "delete":
			message = fmt.Sprintf("Metadata '%s' deleted successfully", key)
		}

		return processComplete{
			success: true,
			message: message,
			context: "save_metadata",
		}
	}
}

func (ui *UIController) handleMetadataSuccessInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.metadataKeys = make([]string, 0, len(ui.decryptedItem.Meta.Map))
		for key := range ui.decryptedItem.Meta.Map {
			ui.metadataKeys = append(ui.metadataKeys, key)
		}
		ui.state = stateMetadataList
		ui.metadataSuccessMsg = ""
		ui.metadataMenu = 0
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) handleMetadataErrorInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.state = stateMetadataList
		ui.metadataErrorMsg = ""
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) metadataSuccessView() string {
	title := successStyle.Render("Metadata Updated Successfully")
	message := ui.metadataSuccessMsg

	controls := "\nControls: Enter to continue, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}

func (ui *UIController) metadataErrorView() string {
	title := errorStyle.Render("Metadata Error")
	message := ui.metadataErrorMsg

	controls := "\nControls: Enter to try again, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}
