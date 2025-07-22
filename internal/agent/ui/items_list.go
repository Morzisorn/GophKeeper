package ui

import (
	"context"
	"fmt"
	"gophkeeper/models"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) handleViewAllItems() (tea.Model, tea.Cmd) {
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

func (ui *UIController) handleViewItemsWithType() (tea.Model, tea.Cmd) {
	ui.state = stateItemTypeSelection
	ui.itemTypeMenu = 0
	return ui, ui.loadItemTypesCmd()
}

func (ui *UIController) loadItemTypesCmd() tea.Cmd {
	return func() tea.Msg {
		itemTypes, err := ui.Item.GetTypesCounts(context.Background(), ui.login)
		if err != nil {
			return errorMsg{
				err:     err,
				context: "load_items_types_counts",
			}
		}
		
		// Убрать типы с нулевым количеством
		for itemType, count := range itemTypes {
			if count == 0 {
				delete(itemTypes, itemType)
			}
		}

		return itemTypesLoaded{
			itemTypes: itemTypes,
		}
	}
}

func (ui *UIController) loadItemsByTypeCmd(itemType string) tea.Cmd {
	return func() tea.Msg {
		items, err := ui.Item.GetItems(context.Background(), ui.login, models.ItemType(itemType))
		if err != nil {
			return errorMsg{
				err:     err,
				context: "load_items_by_type",
			}
		}

		return itemsByTypeLoaded{
			items:    items,
			itemType: itemType,
		}
	}
}

type itemTypesLoaded struct {
	itemTypes map[string]int32
}

type itemsByTypeLoaded struct {
	items    []models.EncryptedItem
	itemType string
}

func (ui *UIController) itemTypeSelectionView() string {
	title := titleStyle.Render("Select Item Type")
	subtitle := "Choose a type to filter items:"

	if ui.addItemErrorMsg != "" {
		errorMsg := errorStyle.Render(ui.addItemErrorMsg)
		return fmt.Sprintf("%s\n\n%s\n\nPress any key to continue", title, errorMsg)
	}

	// Создать слайс для сортированного отображения (сохраняем порядок появления)
	types := make([]string, 0, len(ui.itemTypes))
	for itemType := range ui.itemTypes {
		types = append(types, itemType)
	}

	menu := ""
	for i, itemType := range types {
		count := ui.itemTypes[itemType]
		prefix := fmt.Sprintf("%d. ", i+1)
		option := fmt.Sprintf("%s (%d)", itemType, count)

		if i == ui.itemTypeMenu {
			menu += selectedStyle.Render(prefix+option) + "\n"
		} else {
			menu += menuStyle.Render(prefix+option) + "\n"
		}
	}

	controls := "\nControls: ↑/↓ to navigate, Enter to select, Esc to go back, q to quit"
	return fmt.Sprintf("%s\n\n%s\n\n%s%s", title, subtitle, menu, controls)
}

func (ui *UIController) handleItemTypeSelectionInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Если есть ошибка, любая клавиша возвращает в меню
	if ui.addItemErrorMsg != "" {
		ui.addItemErrorMsg = ""
		ui.state = stateMenuLoggedIn
		return ui, nil
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateMenuLoggedIn
		return ui, nil
	case "up", "k":
		if ui.itemTypeMenu > 0 {
			ui.itemTypeMenu--
		}
	case "down", "j":
		if ui.itemTypeMenu < ui.maxItemTypes {
			ui.itemTypeMenu++
		}
	case "enter":
		return ui.handleViewItemsByType()
	default:
		// Обработка числовых клавиш
		if num := msg.String(); num >= "1" && num <= "9" {
			index := int(num[0] - '1')
			if index <= ui.maxItemTypes {
				ui.itemTypeMenu = index
				return ui.handleViewItemsByType()
			}
		}
	}
	return ui, nil
}

func (ui *UIController) getSelectedType() string {
	// Создать слайс для получения типа по индексу (в том же порядке, что и в представлении)
	types := make([]string, 0, len(ui.itemTypes))
	for itemType := range ui.itemTypes {
		types = append(types, itemType)
	}

	if ui.itemTypeMenu < len(types) {
		return types[ui.itemTypeMenu]
	}
	return ""
}

func (ui *UIController) handleViewItemsByType() (tea.Model, tea.Cmd) {
	ui.state = stateViewItemsByType
	ui.selectedType = ui.getSelectedType()
	return ui, ui.loadItemsByTypeCmd(ui.selectedType)
}

func (ui *UIController) viewItemsByTypeView() string {
	title := titleStyle.Render(fmt.Sprintf("Items with type: %s", ui.selectedType))

	if ui.addItemErrorMsg != "" {
		errorMsg := errorStyle.Render(ui.addItemErrorMsg)
		return fmt.Sprintf("%s\n\n%s\n\nPress any key to continue", title, errorMsg)
	}

	if len(ui.items) == 0 {
		return fmt.Sprintf("%s\n\nNo items found with type '%s'.\n\nPress Esc to go back to type selection, q to quit",
			title, ui.selectedType)
	}

	// Отображение списка айтемов
	menu := ""
	for i, item := range ui.items {
		prefix := fmt.Sprintf("%d. ", i+1)
		itemName := item.Name // или как у вас называется поле с именем
		if itemName == "" {
			itemName = "Unnamed Item"
		}

		if i == ui.currentItem {
			menu += selectedStyle.Render(prefix+itemName) + "\n"
		} else {
			menu += menuStyle.Render(prefix+itemName) + "\n"
		}
	}

	controls := "\nControls: ↑/↓ to navigate, Enter to select, Esc to go back to type selection, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, menu, controls)
}

func (ui *UIController) handleViewItemsByTypeInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Если есть ошибка, вернуться к выбору типа
	if ui.addItemErrorMsg != "" {
		ui.addItemErrorMsg = ""
		ui.state = stateItemTypeSelection
		return ui, nil
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateItemTypeSelection
		return ui, nil
	case "up", "k":
		if ui.currentItem > 0 {
			ui.currentItem--
		}
	case "down", "j":
		if ui.currentItem < ui.maxItems {
			ui.currentItem++
		}
	case "enter":
		if len(ui.items) > 0 {
			ui.selectedItem = &ui.items[ui.currentItem]
			return ui.handleViewItemDetails()
		}
	default:
		// Обработка числовых клавиш для быстрого выбора
		if num := msg.String(); num >= "1" && num <= "9" {
			index := int(num[0] - '1')
			if index <= ui.maxItems {
				ui.currentItem = index
				ui.selectedItem = &ui.items[ui.currentItem]
				return ui.handleViewItemDetails()
			}
		}
	}
	return ui, nil
}
