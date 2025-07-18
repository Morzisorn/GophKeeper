package ui

import (
	"context"
	"fmt"
	"gophkeeper/models"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) handleViewItems() (*UIController, tea.Cmd) {
	ui.state = stateItemsList
	ui.messages.Set("info", "Loading items...")
	ui.currentItem = 0
	return ui, ui.loadItemsCmd()
}

type itemsLoaded struct {
	items []models.EncryptedItem
}

func (ui *UIController) loadItemsCmd() tea.Cmd {
	return func() tea.Msg {
		items, err := ui.Item.GetItems(context.Background(), ui.login, models.ItemTypeUNSPECIFIED)
		if err != nil {
			return errorMsg{
				err:     err,
				context: "load_items",
			}
		}

		return itemsLoaded{
			items: items,
		}
	}
}

func (ui *UIController) handleItemsListInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateMenuLoggedIn
		ui.loggedInMenu = 0
		return ui, nil
	case "up", "k":
		if ui.currentItem > 0 {
			ui.currentItem--
		}
	case "down", "j":
		if len(ui.items) > 0 && ui.currentItem < len(ui.items)-1 {
			ui.currentItem++
		}
	case "enter":
		if len(ui.items) > 0 && ui.currentItem < len(ui.items) {
			return ui.handleViewItemDetails()
		}
	case "r":
		ui.state = stateProcessing
		return ui, ui.loadItemsCmd()
	}
	return ui, nil
}

func (ui *UIController) itemsListView() string {
	title := titleStyle.Render("Your Items")

	if len(ui.items) == 0 {
		controls := "\nControls: r to refresh, Esc to go back, q to quit"
		return fmt.Sprintf("%s\n\nNo items found.%s", title, controls)
	}

	itemsList := ""
	for i, item := range ui.items {
		itemText := fmt.Sprintf("%s (%s)", item.Name, item.Type)
		if i == ui.currentItem {
			itemsList += selectedStyle.Render("→ "+itemText) + "\n"
		} else {
			itemsList += menuStyle.Render("  "+itemText) + "\n"
		}
	}

	controls := "\nControls: ↑/↓ to navigate, Enter to view details, r to refresh, Esc to go back"
	return fmt.Sprintf("%s\n\n%s%s", title, itemsList, controls)
}
