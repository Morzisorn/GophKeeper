package ui

import (
	"context"
	"fmt"
	"gophkeeper/models"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) handleAddItem() (*UIController, tea.Cmd) {
	ui.state = stateAddItem
	ui.input = ""
	ui.itemTypeMenu = 0
	ui.maxItemTypes = 3
	ui.newItem = models.Item{UserLogin: ui.login} 
	ui.addItemErrorMsg = ""                      
	ui.addItemSuccessMsg = ""                    
	return ui, nil
}

func (ui *UIController) handleAddItemSuccessInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.state = stateMenuLoggedIn
		ui.addItemSuccessMsg = ""
		ui.loggedInMenu = 0
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) handleAddItemErrorInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.state = stateMenuLoggedIn
		ui.addItemErrorMsg = ""
		ui.newItem = models.Item{}
		ui.loggedInMenu = 0
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) addItemErrorView() string {
	title := errorStyle.Render("Add Item Error")
	message := ui.addItemErrorMsg

	controls := "\nControls: Enter to return to main menu, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}

func (ui *UIController) addItemSuccessView() string {
	title := successStyle.Render("Item Added Successfully")
	message := ui.addItemSuccessMsg

	controls := "\nControls: Enter to return to main menu, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
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
	case stateAddCardExpiry:
		return ui.handleAddCardExpiryInput(msg)
	case stateAddCardCVV:
		return ui.handleAddCardCVVInput(msg)
	case stateAddCardHolder:
		return ui.handleAddCardHolderInput(msg)
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

		switch ui.newItem.Type {
		case models.ItemTypeCREDENTIALS:
			ui.newItem.Data.(*models.Credentials).Login = data
			ui.input = ""
			ui.state = stateAddCredentialPassword
		case models.ItemTypeTEXT:
			ui.newItem.Data.(*models.Text).Content = data
			return ui, ui.addItemCmd()
		case models.ItemTypeCARD:
			if err := validateCardNumber(data); err != nil {
				ui.messages.Set("error", err.Error())
				return ui, nil
			}
			ui.newItem.Data.(*models.Card).Number = data
			ui.input = ""
			ui.state = stateAddCardExpiry
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

func (ui *UIController) addItemPasswordView() string {
	title := titleStyle.Render("Add Credentials - Enter Password")
	input := inputStyle.Render(ui.input + "█")

	controls := "\nControls: Esc to go back, Enter to save"
	return fmt.Sprintf("%s\n\nName: %s\nLogin: %s\nPassword: %s%s",
		title, ui.newItem.Name, ui.newItem.Data.(*models.Credentials).Login, input, controls)
}

func (ui *UIController) addItemDataView() string {
	var prompt string
	var hint string
	
	switch ui.newItem.Type {
	case models.ItemTypeCREDENTIALS:
		prompt = "Login"
	case models.ItemTypeTEXT:
		prompt = "Text Content"
	case models.ItemTypeBINARY:
		prompt = "Binary Data"
	case models.ItemTypeCARD:
		prompt = "Card Number"
		hint = " (16 or 18 digits)"
	default:
		prompt = "Data"
	}

	title := titleStyle.Render(fmt.Sprintf("Add %s - Enter %s%s", ui.newItem.Type, prompt, hint))
	input := inputStyle.Render(ui.input + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to continue"
	return fmt.Sprintf("%s\n\nName: %s\n%s: %s%s%s",
		title, ui.newItem.Name, prompt, input, errorMsg, controls)
}

func (ui *UIController) handleAddCardExpiryInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return ui, tea.Quit
	case "esc":
		ui.state = stateAddItemData
		ui.input = ui.newItem.Data.(*models.Card).Number
		return ui, nil
	case "enter":
		expiry := strings.TrimSpace(ui.input)
		if expiry == "" {
			ui.messages.Set("error", "Expiry date cannot be empty")
			return ui, nil
		}
		
		if err := validateExpiry(expiry); err != nil {
			ui.messages.Set("error", err.Error())
			return ui, nil
		}
		
		ui.newItem.Data.(*models.Card).ExpiryDate = expiry
		ui.input = ""
		ui.state = stateAddCardCVV
		ui.messages.Clear("error")
		return ui, nil
	case "backspace":
		if len(ui.input) > 0 {
			ui.input = ui.input[:len(ui.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			if len(ui.input) == 2 && msg.String() != "/" {
				ui.input += "/"
			}
			ui.input += msg.String()
		}
	}
	return ui, nil
}

func (ui *UIController) handleAddCardCVVInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return ui, tea.Quit
	case "esc":
		ui.state = stateAddCardExpiry
		ui.input = ui.newItem.Data.(*models.Card).ExpiryDate
		return ui, nil
	case "enter":
		cvv := strings.TrimSpace(ui.input)
		if cvv == "" {
			ui.messages.Set("error", "CVV cannot be empty")
			return ui, nil
		}
		
		if err := validateCVV(cvv); err != nil {
			ui.messages.Set("error", err.Error())
			return ui, nil
		}
		
		ui.newItem.Data.(*models.Card).SecurityCode = cvv
		ui.input = ""
		ui.state = stateAddCardHolder
		ui.messages.Clear("error")
		return ui, nil
	case "backspace":
		if len(ui.input) > 0 {
			ui.input = ui.input[:len(ui.input)-1]
		}
	default:
		if len(msg.String()) == 1 && regexp.MustCompile(`\d`).MatchString(msg.String()) && len(ui.input) < 3 {
			ui.input += msg.String()
		}
	}
	return ui, nil
}

func (ui *UIController) handleAddCardHolderInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return ui, tea.Quit
	case "esc":
		ui.state = stateAddCardCVV
		ui.input = ui.newItem.Data.(*models.Card).SecurityCode
		return ui, nil
	case "enter":
		holder := strings.TrimSpace(ui.input)
		if holder == "" {
			ui.messages.Set("error", "Cardholder name cannot be empty")
			return ui, nil
		}
		ui.newItem.Data.(*models.Card).CardholderName = holder
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

func (ui *UIController) addCardExpiryView() string {
	title := titleStyle.Render("Add Card - Enter Expiry Date")
	input := inputStyle.Render(ui.input + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to continue"
	hint := "\nFormat: MM/YY (e.g., 12/25)"
	return fmt.Sprintf("%s\n\nName: %s\nNumber: %s\nExpiry: %s%s%s%s",
		title, ui.newItem.Name, ui.newItem.Data.(*models.Card).Number, input, hint, errorMsg, controls)
}

func (ui *UIController) addCardCVVView() string {
	title := titleStyle.Render("Add Card - Enter CVV")
	hiddenCVV := strings.Repeat("*", len(ui.input))
	input := inputStyle.Render(hiddenCVV + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to continue"
	hint := "\n3 digits only"
	return fmt.Sprintf("%s\n\nName: %s\nNumber: %s\nExpiry: %s\nCVV: %s%s%s%s",
		title, ui.newItem.Name, ui.newItem.Data.(*models.Card).Number, 
		ui.newItem.Data.(*models.Card).ExpiryDate, input, hint, errorMsg, controls)
}

func (ui *UIController) addCardHolderView() string {
	title := titleStyle.Render("Add Card - Enter Cardholder Name")
	input := inputStyle.Render(ui.input + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to save"
	return fmt.Sprintf("%s\n\nName: %s\nNumber: %s\nExpiry: %s\nCVV: %s\nCardholder: %s%s%s",
		title, ui.newItem.Name, ui.newItem.Data.(*models.Card).Number,
		ui.newItem.Data.(*models.Card).ExpiryDate, strings.Repeat("*", len(ui.newItem.Data.(*models.Card).SecurityCode)),
		input, errorMsg, controls)
}

func validateCardNumber(number string) error {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(number, " ", ""), "-", "")

	if !regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
		return fmt.Errorf("card number must contain only digits")
	}

	if len(cleaned) != 16 && len(cleaned) != 18 {
		return fmt.Errorf("card number must be 16 or 18 digits")
	}

	return nil
}

func validateExpiry(expiry string) error {
	matched, _ := regexp.MatchString(`^\d{2}/\d{2}$`, expiry)
	if !matched {
		return fmt.Errorf("expiry must be in MM/YY format")
	}

	parts := strings.Split(expiry, "/")
	month := parts[0]

	if month < "01" || month > "12" {
		return fmt.Errorf("month must be between 01 and 12")
	}

	return nil
}

func validateCVV(cvv string) error {
	if !regexp.MustCompile(`^\d{3}$`).MatchString(cvv) {
		return fmt.Errorf("CVV must be exactly 3 digits")
	}

	return nil
}
