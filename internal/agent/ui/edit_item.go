package ui

import (
	"context"
	"fmt"
	"gophkeeper/models"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) startEditItem() (*UIController, tea.Cmd) {
	if ui.decryptedItem == nil {
		return ui, nil
	}

	ui.editingItem = &models.Item{
		ID:        ui.decryptedItem.ID,
		UserLogin: ui.decryptedItem.UserLogin,
		Name:      ui.decryptedItem.Name,
		Type:      ui.decryptedItem.Type,
		Meta:      ui.decryptedItem.Meta,
		CreatedAt: ui.decryptedItem.CreatedAt,
		UpdatedAt: ui.decryptedItem.UpdatedAt,
	}

	switch data := ui.decryptedItem.Data.(type) {
	case *models.Credentials:
		ui.editingItem.Data = &models.Credentials{
			Login:    data.Login,
			Password: data.Password,
		}
	case *models.Text:
		ui.editingItem.Data = &models.Text{
			Content: data.Content,
		}
	case *models.Card:
		ui.editingItem.Data = &models.Card{
			Number:         data.Number,
			ExpiryDate:     data.ExpiryDate,
			SecurityCode:   data.SecurityCode,
			CardholderName: data.CardholderName,
		}
	case *models.Binary:
		ui.editingItem.Data = &models.Binary{
			Content: data.Content,
		}
	}
	ui.state = stateEditItemName
	ui.input = ui.editingItem.Name
	ui.editStep = 0
	return ui, nil
}

func (ui *UIController) handleEditItemNameInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.editingItem = nil
		ui.state = stateItemDetails
		ui.input = ""
		return ui, nil
	case "enter":
		name := strings.TrimSpace(ui.input)
		if name == "" {
			return ui, nil
		}
		ui.editingItem.Name = name
		ui.input = ""

		switch ui.editingItem.Type {
		case models.ItemTypeCREDENTIALS:
			ui.state = stateEditCredentialLogin
			ui.input = ui.editingItem.Data.(*models.Credentials).Login
		case models.ItemTypeTEXT:
			ui.state = stateEditTextContent
			ui.input = ui.editingItem.Data.(*models.Text).Content
		case models.ItemTypeCARD:
			ui.state = stateEditCardNumber
			ui.input = ui.editingItem.Data.(*models.Card).Number
		case models.ItemTypeBINARY:
			ui.state = stateEditBinaryData
			ui.input = string(ui.editingItem.Data.(*models.Binary).Content)
		}
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

func (ui *UIController) handleEditCredentialLoginInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateEditItemName
		ui.input = ui.editingItem.Name
		return ui, nil
	case "enter":
		login := strings.TrimSpace(ui.input)
		ui.editingItem.Data.(*models.Credentials).Login = login
		ui.state = stateEditCredentialPassword
		ui.input = ui.editingItem.Data.(*models.Credentials).Password
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

func (ui *UIController) handleEditCredentialPasswordInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateEditCredentialLogin
		ui.input = ui.editingItem.Data.(*models.Credentials).Login
		return ui, nil
	case "enter":
		password := strings.TrimSpace(ui.input)
		ui.editingItem.Data.(*models.Credentials).Password = password
		ui.input = ""
		return ui, ui.saveEditedItemCmd()
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

func (ui *UIController) handleEditItemDataInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateEditItemName
		ui.input = ui.editingItem.Name
		return ui, nil
	case "enter":
		data := strings.TrimSpace(ui.input)

		switch ui.editingItem.Type {
		case models.ItemTypeTEXT:
			ui.editingItem.Data.(*models.Text).Content = data
		case models.ItemTypeBINARY:
			ui.editingItem.Data.(*models.Binary).Content = []byte(data)
		}

		ui.input = ""
		return ui, ui.saveEditedItemCmd()
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

func (ui *UIController) saveEditedItemCmd() tea.Cmd {
	return func() tea.Msg {
		ui.editingItem.UpdatedAt = time.Now()

		err := ui.Item.EditItem(context.Background(), ui.editingItem)
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Edit error: %v", err),
				context: "edit_item",
			}
		}

		return processComplete{
			success: true,
			message: "Item updated successfully",
			context: "edit_item",
		}
	}
}

func (ui *UIController) handleEditSuccessInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.state = stateProcessing
		ui.editSuccessMsg = ""
		ui.editingItem = nil
		return ui, ui.loadItemsCmd()
	}
	return ui, nil
}

func (ui *UIController) handleEditErrorInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		ui.state = stateItemDetails
		ui.editErrorMsg = ""
		ui.editingItem = nil
		return ui, nil
	case "esc":
		ui.state = stateItemDetails
		ui.editErrorMsg = ""
		ui.editingItem = nil
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) editItemNameView() string {
	title := titleStyle.Render("Edit Item - Name")
	input := inputStyle.Render(ui.input + "█")

	controls := "\nControls: Esc to cancel, Enter to continue"
	return fmt.Sprintf("%s\n\nName: %s%s", title, input, controls)
}

func (ui *UIController) editCredentialLoginView() string {
	title := titleStyle.Render("Edit Credentials - Login")
	input := inputStyle.Render(ui.input + "█")

	controls := "\nControls: Esc to go back, Enter to continue"
	return fmt.Sprintf("%s\n\nName: %s\nLogin: %s%s",
		title, ui.editingItem.Name, input, controls)
}

func (ui *UIController) editCredentialPasswordView() string {
	title := titleStyle.Render("Edit Credentials - Password")
	hiddenPassword := strings.Repeat("*", len(ui.input))
	input := inputStyle.Render(hiddenPassword + "█")

	controls := "\nControls: Esc to go back, Enter to save"
	return fmt.Sprintf("%s\n\nName: %s\nLogin: %s\nPassword: %s%s",
		title, ui.editingItem.Name,
		ui.editingItem.Data.(*models.Credentials).Login, input, controls)
}

func (ui *UIController) editItemDataView() string {
	var dataType string
	switch ui.editingItem.Type {
	case models.ItemTypeTEXT:
		dataType = "Text Content"
	case models.ItemTypeBINARY:
		dataType = "Binary Data"
	default:
		dataType = "Data"
	}

	title := titleStyle.Render(fmt.Sprintf("Edit %s - %s", ui.editingItem.Type, dataType))
	input := inputStyle.Render(ui.input + "█")

	controls := "\nControls: Esc to go back, Enter to save"
	return fmt.Sprintf("%s\n\nName: %s\n%s: %s%s",
		title, ui.editingItem.Name, dataType, input, controls)
}

func (ui *UIController) editSuccessView() string {
	title := successStyle.Render("Item Updated Successfully")
	message := ui.editSuccessMsg

	controls := "\nControls: Enter to return to items list, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}

func (ui *UIController) editErrorView() string {
	title := errorStyle.Render("Edit Error")
	message := ui.editErrorMsg

	controls := "\nControls: Enter to return to item details, Esc to cancel, q to quit"
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}

func (ui *UIController) handleEditTextContentInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateEditItemName
		ui.input = ui.editingItem.Name
		return ui, nil
	case "enter":
		content := strings.TrimSpace(ui.input)
		if content == "" {
			return ui, nil
		}

		ui.editingItem.Data.(*models.Text).Content = content
		ui.input = ""
		return ui, ui.saveEditedItemCmd()
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

func (ui *UIController) editTextContentView() string {
	title := titleStyle.Render("Edit Text - Content")
	input := inputStyle.Render(ui.input + "█")

	controls := "\nControls: Esc to go back, Enter to save"
	return fmt.Sprintf("%s\n\nName: %s\nContent: %s%s",
		title, ui.editingItem.Name, input, controls)
}

func (ui *UIController) handleEditBinaryDataInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateEditItemName
		ui.input = ui.editingItem.Name
		return ui, nil
	case "enter":
		data := strings.TrimSpace(ui.input)
		if data == "" {
			return ui, nil
		}

		ui.editingItem.Data.(*models.Binary).Content = []byte(data)
		ui.input = ""
		return ui, ui.saveEditedItemCmd()
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

func (ui *UIController) editBinaryDataView() string {
	title := titleStyle.Render("Edit Binary - Data")
	input := inputStyle.Render(ui.input + "█")

	controls := "\nControls: Esc to go back, Enter to save"
	return fmt.Sprintf("%s\n\nName: %s\nData: %s%s",
		title, ui.editingItem.Name, input, controls)
}

func (ui *UIController) handleEditCardNumberInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateEditItemName
		ui.input = ui.editingItem.Name
		return ui, nil
	case "enter":
		number := strings.TrimSpace(ui.input)
		if number == "" {
			ui.messages.Set("error", "Card number cannot be empty")
			return ui, nil
		}

		if err := validateCardNumber(number); err != nil {
			ui.messages.Set("error", err.Error())
			return ui, nil
		}

		ui.editingItem.Data.(*models.Card).Number = number
		ui.state = stateEditCardExpiry
		ui.input = ui.editingItem.Data.(*models.Card).ExpiryDate
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

func (ui *UIController) handleEditCardExpiryInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateEditCardNumber
		ui.input = ui.editingItem.Data.(*models.Card).Number
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

		ui.editingItem.Data.(*models.Card).ExpiryDate = expiry
		ui.state = stateEditCardCVV
		ui.input = ui.editingItem.Data.(*models.Card).SecurityCode
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

func (ui *UIController) handleEditCardCVVInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateEditCardExpiry
		ui.input = ui.editingItem.Data.(*models.Card).ExpiryDate
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

		ui.editingItem.Data.(*models.Card).SecurityCode = cvv
		ui.state = stateEditCardHolder
		ui.input = ui.editingItem.Data.(*models.Card).CardholderName
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

func (ui *UIController) handleEditCardHolderInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "esc":
		ui.state = stateEditCardCVV
		ui.input = ui.editingItem.Data.(*models.Card).SecurityCode
		return ui, nil
	case "enter":
		holder := strings.TrimSpace(ui.input)
		if holder == "" {
			return ui, nil
		}
		ui.editingItem.Data.(*models.Card).CardholderName = holder
		ui.input = ""
		return ui, ui.saveEditedItemCmd()
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

func (ui *UIController) editCardNumberView() string {
	title := titleStyle.Render("Edit Card - Number")
	input := inputStyle.Render(ui.input + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to continue"
	hint := "\n16 or 18 digits only"
	return fmt.Sprintf("%s\n\nName: %s\nNumber: %s%s%s%s",
		title, ui.editingItem.Name, input, hint, errorMsg, controls)
}

func (ui *UIController) editCardExpiryView() string {
	title := titleStyle.Render("Edit Card - Expiry Date")
	input := inputStyle.Render(ui.input + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to continue"
	hint := "\nFormat: MM/YY (e.g., 12/25)"
	return fmt.Sprintf("%s\n\nName: %s\nNumber: %s\nExpiry: %s%s%s%s",
		title, ui.editingItem.Name, ui.editingItem.Data.(*models.Card).Number, input, hint, errorMsg, controls)
}

func (ui *UIController) editCardCVVView() string {
	title := titleStyle.Render("Edit Card - CVV")
	hiddenCVV := strings.Repeat("*", len(ui.input))
	input := inputStyle.Render(hiddenCVV + "█")

	errorMsg := ""
	if err := ui.messages.Get("error"); err != "" {
		errorMsg = "\n" + errorStyle.Render(err)
	}

	controls := "\nControls: Esc to go back, Enter to continue"
	hint := "\n3 digits only"
	return fmt.Sprintf("%s\n\nName: %s\nNumber: %s\nExpiry: %s\nCVV: %s%s%s%s",
		title, ui.editingItem.Name, ui.editingItem.Data.(*models.Card).Number,
		ui.editingItem.Data.(*models.Card).ExpiryDate, input, hint, errorMsg, controls)
}

func (ui *UIController) editCardHolderView() string {
	title := titleStyle.Render("Edit Card - Cardholder Name")
	input := inputStyle.Render(ui.input + "█")

	controls := "\nControls: Esc to go back, Enter to save"
	return fmt.Sprintf("%s\n\nName: %s\nNumber: %s\nExpiry: %s\nCVV: %s\nCardholder: %s%s",
		title, ui.editingItem.Name, ui.editingItem.Data.(*models.Card).Number,
		ui.editingItem.Data.(*models.Card).ExpiryDate, strings.Repeat("*", len(ui.editingItem.Data.(*models.Card).SecurityCode)),
		input, controls)
}
