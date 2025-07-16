package ui

import (
	"context"
	"fmt"
	"gophkeeper/models"
	"strings"

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
		if ui.loggedInMenu < ui.maxLoggedInMenu {
			ui.loggedInMenu++
		}
	case "1":
		ui.loggedInMenu = 0
		return ui.handleViewItems()
	case "2":
		ui.loggedInMenu = 1
		return ui.handleAddItem()
	// case "3":
	// 	ui.loggedInMenu = 2
	// 	return ui.handleLogout()
	case "enter":
		switch ui.loggedInMenu {
		case 0:
			return ui.handleViewItems()
		case 1:
			return ui.handleAddItem()
		// case 2:
		// 	return ui.handleLogout(), nil
		}
	}
	return ui, nil
}

func (ui *UIController) handleViewItems() (*UIController, tea.Cmd) {
	ui.state = stateItemsList
	ui.messages.Set("info", "Loading items...")
	ui.currentItem = 0
	return ui, ui.loadItemsCmd()
}

func (ui *UIController) handleAddItem() (*UIController, tea.Cmd) {
	ui.state = stateAddItem
	ui.input = ""
	ui.itemTypeMenu = 0
	ui.maxItemTypes = 3        // LOGIN, CARD, SECURE_NOTE, BINARY
	ui.newItem = models.Item{UserLogin: ui.login} // Очищаем данные
	ui.messages.ClearAll()
	return ui, nil
}

// func (ui *UIController) handleLogout() (*UIController, tea.Cmd) {
// 	// Выход из системы
// 	ui.isAuthenticated = false
// 	ui.state = stateMenuLoggedOut
// 	return ui, ui.logoutCmd()
// }

type itemsLoaded struct {
	items []models.Item
}

func (ui *UIController) loadItemsCmd() tea.Cmd {
	return func() tea.Msg {
		// Вызываем сервис для получения элементов
		items, err := ui.Item.GetItems(context.Background(), ui.login, "")
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
		// Возвращаемся в главное меню
		ui.state = stateMenuLoggedIn
		ui.loggedInMenu = 0
		ui.messages.ClearAll()
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
			// Здесь можно добавить просмотр деталей элемента
			return ui.handleViewItemDetails()
		}
	case "d":
		if len(ui.items) > 0 {
			// Удаление элемента
			return ui.handleDeleteItem()
		}
	}
	return ui, nil
}

func (ui *UIController) itemsListView() string {
	title := titleStyle.Render("Your Items")

	// Если есть сообщение о загрузке
	if info := ui.messages.Get("info"); info != "" {
		return fmt.Sprintf("%s\n\n%s", title, info)
	}

	// Если нет элементов
	if len(ui.items) == 0 {
		controls := "\nControls: Esc to go back, q to quit"
		return fmt.Sprintf("%s\n\nNo items found.%s", title, controls)
	}

	// Отображаем список элементов
	itemsList := ""
	for i, item := range ui.items {
		itemText := fmt.Sprintf("%s (%s)", item.Name, item.Type)
		if i == ui.currentItem {
			itemsList += selectedStyle.Render("→ "+itemText) + "\n"
		} else {
			itemsList += menuStyle.Render("  "+itemText) + "\n"
		}
	}

	controls := "\nControls: ↑/↓ to navigate, Enter to view details, d to delete, Esc to go back"
	return fmt.Sprintf("%s\n\n%s%s", title, itemsList, controls)
}

func (ui *UIController) handleViewItemDetails() (*UIController, tea.Cmd) {
	if ui.currentItem < len(ui.items) {
		selectedItem := ui.items[ui.currentItem]
		ui.messages.Set("item_details", fmt.Sprintf("Name: %s\nType: %s\nData: %s\nMeta: %s",
			selectedItem.Name, selectedItem.Type, selectedItem.Data, selectedItem.Meta))
	}
	return ui, nil
}

func (ui *UIController) handleDeleteItem() (*UIController, tea.Cmd) {
	if ui.currentItem < len(ui.items) {
		selectedItem := ui.items[ui.currentItem]
		ui.state = stateProcessing
		return ui, ui.deleteItemCmd(selectedItem.ID)
	}
	return ui, nil
}

func (ui *UIController) deleteItemCmd(itemID string) tea.Cmd {
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

func (ui *UIController) handleAddItemInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch ui.state {
	case stateAddItem:
		return ui.handleItemTypeSelection(msg)
	case stateAddItemName:
		return ui.handleItemNameInput(msg)
	case stateAddItemData:
		return ui.handleItemDataInput(msg)
	case stateAddCredentialPassword:
		return ui.handleCredentialPasswordInput(msg)
	}
	return ui, nil
}

func (ui *UIController) handleItemTypeSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateMenuLoggedIn
		ui.messages.ClearAll()
		return ui, nil
	case "up", "k":
		if ui.itemTypeMenu > 0 {
			ui.itemTypeMenu--
		}
	case "down", "j":
		if ui.itemTypeMenu < ui.maxItemTypes {
			ui.itemTypeMenu++
		}
	case "1":
		ui.itemTypeMenu = 0
		return ui.selectItemType(models.ItemTypeCREDENTIALS), nil
	case "2":
		ui.itemTypeMenu = 1
		return ui.selectItemType(models.ItemTypeTEXT), nil
	case "3":
		ui.itemTypeMenu = 2
		return ui.selectItemType(models.ItemTypeBINARY), nil
	case "4":
		ui.itemTypeMenu = 3
		return ui.selectItemType(models.ItemTypeCARD), nil
	case "enter":
		types := models.ItemTypes
		return ui.selectItemType(types[ui.itemTypeMenu]), nil
	}
	return ui, nil
}

func (ui *UIController) selectItemType(itemType models.ItemType) *UIController {
	ui.newItem.Type = itemType
	
	// Инициализируем структуру данных в зависимости от типа
	switch itemType {
	case models.ItemTypeCREDENTIALS:
		ui.newItem.Data = &models.Credentials{}
	case models.ItemTypeCARD:
		ui.newItem.Data = &models.Card{}
	case models.ItemTypeTEXT:
		ui.newItem.Data = &models.Text{}
	case models.ItemTypeBINARY:
		ui.newItem.Data = &models.Binary{}
	}
	
	ui.state = stateAddItemName
	ui.input = ""
	return ui
}

func (ui *UIController) handleItemNameInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return ui, tea.Quit
	case "esc":
		ui.state = stateAddItem
		ui.input = ""
		return ui, nil
	case "enter":
		ui.newItem.Name = strings.TrimSpace(ui.input)
		if ui.newItem.Name == "" {
			ui.messages.Set("error", "Name cannot be empty")
			return ui, nil
		}
		ui.input = ""
		ui.state = stateAddItemData
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

func (ui *UIController) handleItemDataInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return ui, tea.Quit
	case "esc":
		ui.state = stateAddItemName
		ui.input = ui.newItem.Name
		return ui, nil
	case "enter":
		data := strings.TrimSpace(ui.input)
		if data == "" {
			ui.messages.Set("error", "Data cannot be empty")
			return ui, nil
		}

		// Сохраняем данные в зависимости от типа
		switch ui.newItem.Type {
		case models.ItemTypeCREDENTIALS:
			ui.newItem.Data.(*models.Credentials).Login = data
			ui.input = ""
			ui.state = stateAddCredentialPassword
		case models.ItemTypeCARD:
			ui.newItem.Data.(*models.Card).Number = data
			ui.input = ""
			// Для карты можно сразу сохранить или перейти к вводу других данных
			return ui, ui.addItemCmd()
		case models.ItemTypeTEXT:
			ui.newItem.Data.(*models.Text).Content = data
			return ui, ui.addItemCmd()
		case models.ItemTypeBINARY:
			ui.newItem.Data.(*models.Binary).Content = []byte(data)
			return ui, ui.addItemCmd()
		}
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

func (ui *UIController) handleCredentialPasswordInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return ui, tea.Quit
	case "esc":
		ui.state = stateAddItem
		ui.input = ui.newItem.Data.(*models.Credentials).Login
		return ui, nil
	case "enter":
		ui.newItem.Data.(*models.Credentials).Password = strings.TrimSpace(ui.input)
		ui.input = ""
		return ui, ui.addItemCmd()
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

func (ui *UIController) addItemCmd() tea.Cmd {
	return func() tea.Msg {
		// Вызываем сервис для сохранения элемента
		err := ui.Item.AddItem(context.Background(), &ui.newItem)
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Save error: %v", err),
				context: "save_item",
			}
		}

		return processComplete{
			success: true,
			message: "Item saved successfully",
			context: "save_item",
		}
	}
}

func (ui *UIController) addItemTypeView() string {
	title := titleStyle.Render("Add New Item - Select Type")

	types := []string{"Credentials", "Text", "Binary", "Credit Card"}

	menu := ""
	for i, itemType := range types {
		prefix := fmt.Sprintf("%d. ", i+1)
		if i == ui.itemTypeMenu {
			menu += selectedStyle.Render(prefix+itemType) + "\n"
		} else {
			menu += menuStyle.Render(prefix+itemType) + "\n"
		}
	}

	controls := "\nControls: ↑/↓ to navigate, Enter to select, Esc to go back"
	return fmt.Sprintf("%s\n\n%s%s", title, menu, controls)
}

func (ui *UIController) addItemNameView() string {
	title := titleStyle.Render(fmt.Sprintf("Add %s - Enter Name", ui.newItem.Type))
	input := inputStyle.Render(ui.input + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to continue"
	return fmt.Sprintf("%s\n\nName: %s%s%s", title, input, errorMsg, controls)
}

func (ui *UIController) addItemDataView() string {
	var prompt string
	switch ui.newItem.Type {
	case models.ItemTypeCREDENTIALS:
		prompt = "Login"
	case models.ItemTypeTEXT:
		prompt = "Text"
	case models.ItemTypeBINARY:
		prompt = "Binary"
	case models.ItemTypeCARD:
		prompt = "Card"
	default:
		prompt = "Unknown"
	}

	title := titleStyle.Render(fmt.Sprintf("Add %s - Enter %s", ui.newItem.Type, prompt))
	input := inputStyle.Render(ui.input + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to continue"
	return fmt.Sprintf("%s\n\nName: %s\n%s: %s%s%s",
		title, ui.newItem.Name, prompt, input, errorMsg, controls)
}

func (ui *UIController) addItemPasswordView() string {
	title := titleStyle.Render("Add Credentials - Enter Password")
	hiddenPassword := strings.Repeat("*", len(ui.input))
	input := inputStyle.Render(hiddenPassword + "█")

	controls := "\nControls: Esc to go back, Enter to save"
	return fmt.Sprintf("%s\n\nName: %s\nLogin: %s\nPassword: %s%s",
		title, ui.newItem.Name, ui.newItem.Data.(*models.Credentials).Login, input, controls)
}
