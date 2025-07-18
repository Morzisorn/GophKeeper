package ui

import (
	"fmt"
	"gophkeeper/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	processComplete struct {
		success bool
		message string
		context string
	}
	errorMsg struct {
		err     error
		context string
	}
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	menuStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			MarginLeft(2)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#EE6FF8")).
			Background(lipgloss.Color("#3C3C3C"))

	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true)
)

func (ui *UIController) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return ui.handleKeyMsg(msg)
	case itemDecrypted:
		ui.decryptedItem = msg.item
		return ui, nil
	case decryptError:
		return ui.handleDecryptError(msg)
	case processComplete:
		return ui.handleProcessComplete(msg)
	case itemsLoaded:
		ui.items = msg.items
		ui.currentItem = 0
		ui.state = stateItemsList
		return ui, nil
	case errorMsg:
		ui.state = stateItemsList
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case ui.state == stateMenuLoggedOut:
		return ui.handleMenuLoggedOutInput(msg)
	case ui.state == stateMenuLoggedIn:
		return ui.handleMenuLoggedInInput(msg)
	case ui.state == stateMasterPassword:
		return ui.handleMasterPasswordInput(msg)
	case ui.state == stateItemsList:
		return ui.handleItemsListInput(msg)
	case ui.state == stateItemDetails:
		return ui.handleItemDetailsInput(msg)
	case ui.state == stateConfirmDelete:
		return ui.handleConfirmDeleteInput(msg)
	case ui.state == stateDecryptError:
		return ui.handleDecryptErrorInput(msg)
	case ui.state == stateDeleteSuccess:
		return ui.handleDeleteSuccessInput(msg)
	case ui.state == stateDeleteError:
		return ui.handleDeleteErrorInput(msg)
	case ui.state == stateEditItemName:
		return ui.handleEditItemNameInput(msg)
	case ui.state == stateEditCredentialLogin:
		return ui.handleEditCredentialLoginInput(msg)
	case ui.state == stateEditCredentialPassword:
		return ui.handleEditCredentialPasswordInput(msg)
	case ui.state == stateEditItemData:
		return ui.handleEditItemDataInput(msg)
	case ui.state == stateEditSuccess:
		return ui.handleEditSuccessInput(msg)
	case ui.state == stateEditError:
		return ui.handleEditErrorInput(msg)
	case ui.state == stateEditTextContent:
		return ui.handleEditTextContentInput(msg)
	case ui.state == stateEditBinaryData:
		return ui.handleEditBinaryDataInput(msg)
	case ui.state == stateAddItemError:
		return ui.handleAddItemErrorInput(msg)
	case ui.state == stateAddItemSuccess:
		return ui.handleAddItemSuccessInput(msg)
	case ui.state.IsAddItemInput():
		return ui.handleAddItemInput(msg)
	case ui.state == stateAddCardExpiry:
		return ui.handleAddCardExpiryInput(msg)
	case ui.state == stateAddCardCVV:
		return ui.handleAddCardCVVInput(msg)
	case ui.state == stateAddCardHolder:
		return ui.handleAddCardHolderInput(msg)
	case ui.state == stateEditCardNumber:
		return ui.handleEditCardNumberInput(msg)
	case ui.state == stateEditCardExpiry:
		return ui.handleEditCardExpiryInput(msg)
	case ui.state == stateEditCardCVV:
		return ui.handleEditCardCVVInput(msg)
	case ui.state == stateEditCardHolder:
		return ui.handleEditCardHolderInput(msg)
	case ui.state.IsLoginInput():
		return ui.handleLoginInput(msg)
	case ui.state.IsPasswordInput():
		return ui.handlePasswordInput(msg)
	case ui.state == stateSuccess || ui.state == stateError:
		return ui.handleResultInput(msg)
	case ui.state == stateProcessing:
		return ui, nil
	case ui.state == stateMetadataList:
		return ui.handleMetadataListInput(msg)
	case ui.state == stateAddMetadataKey:
		return ui.handleAddMetadataKeyInput(msg)
	case ui.state == stateAddMetadataValue:
		return ui.handleAddMetadataValueInput(msg)
	case ui.state == stateEditMetadataValue:
		return ui.handleEditMetadataValueInput(msg)
	case ui.state == stateConfirmDeleteMetadata:
		return ui.handleConfirmDeleteMetadataInput(msg)
	case ui.state == stateMetadataSuccess:
		return ui.handleMetadataSuccessInput(msg)
	case ui.state == stateMetadataError:
		return ui.handleMetadataErrorInput(msg)
	case ui.state == stateConfirmLogout:
		return ui.handleConfirmLogoutInput(msg)
	case ui.state == stateLogoutSuccess:
		return ui.handleLogoutSuccessInput(msg)
	case ui.state == stateLogoutError:
		return ui.handleLogoutErrorInput(msg)
	}
	return ui, nil
}

func (ui *UIController) View() string {
	debug := fmt.Sprintf("Current state: %d\n", int(ui.state))
	switch {
	case ui.state == stateMenuLoggedOut:
		return ui.menuLoggedOutView()
	case ui.state == stateMenuLoggedIn:
		return ui.menuLoggedInView()
	case ui.state == stateItemsList:
		return ui.itemsListView()
	case ui.state == stateItemDetails:
		return ui.itemDetailsView()
	case ui.state == stateConfirmDelete:
		return ui.confirmDeleteView()
	case ui.state == stateDecryptError:
		return ui.decryptErrorView()
	case ui.state == stateDeleteSuccess:
		return ui.deleteSuccessView()
	case ui.state == stateDeleteError:
		return ui.deleteErrorView()
	case ui.state == stateEditItemName:
		return ui.editItemNameView()
	case ui.state == stateEditCredentialLogin:
		return ui.editCredentialLoginView()
	case ui.state == stateEditCredentialPassword:
		return ui.editCredentialPasswordView()
	case ui.state == stateEditItemData:
		return ui.editItemDataView()
	case ui.state == stateEditSuccess:
		return ui.editSuccessView()
	case ui.state == stateEditError:
		return ui.editErrorView()
	case ui.state == stateEditTextContent:
		return ui.editTextContentView()
	case ui.state == stateEditBinaryData:
		return ui.editBinaryDataView()
	case ui.state == stateAddItemError:
		return ui.addItemErrorView()
	case ui.state == stateAddItemSuccess:
		return ui.addItemSuccessView()
	case ui.state.IsLoginInput():
		return ui.loginInputView()
	case ui.state.IsPasswordInput():
		return ui.passwordInputView()
	case ui.state == stateMasterPassword:
		return ui.masterPasswordInputView()
	case ui.state == stateProcessing:
		return ui.processingView()
	case ui.state == stateSuccess:
		return ui.successView()
	case ui.state == stateError:
		return ui.errorView()
	case ui.state == stateAddItem:
		return ui.addItemTypeView()
	case ui.state == stateAddItemName:
		return ui.addItemNameView()
	case ui.state == stateAddItemData:
		return ui.addItemDataView()
	case ui.state == stateAddCredentialPassword:
		return ui.addItemPasswordView()
	case ui.state == stateAddCardExpiry:
		return ui.addCardExpiryView()
	case ui.state == stateAddCardCVV:
		return ui.addCardCVVView()
	case ui.state == stateAddCardHolder:
		return ui.addCardHolderView()
	case ui.state == stateEditCardNumber:
		return ui.editCardNumberView()
	case ui.state == stateEditCardExpiry:
		return ui.editCardExpiryView()
	case ui.state == stateEditCardCVV:
		return ui.editCardCVVView()
	case ui.state == stateEditCardHolder:
		return ui.editCardHolderView()
	case ui.state == stateMetadataList:
		return ui.metadataListView()
	case ui.state == stateAddMetadataKey:
		return ui.addMetadataKeyView()
	case ui.state == stateAddMetadataValue:
		return ui.addMetadataValueView()
	case ui.state == stateEditMetadataValue:
		return ui.editMetadataValueView()
	case ui.state == stateConfirmDeleteMetadata:
		return ui.confirmDeleteMetadataView()
	case ui.state == stateMetadataSuccess:
		return ui.metadataSuccessView()
	case ui.state == stateMetadataError:
		return ui.metadataErrorView()
	case ui.state == stateConfirmLogout:
		return ui.confirmLogoutView()
	case ui.state == stateLogoutSuccess:
		return ui.logoutSuccessView()
	case ui.state == stateLogoutError:
		return ui.logoutErrorView()
	}
	return "View error:" + debug
}

func (ui *UIController) successView() string {
	title := successStyle.Render("Success!")
	message := ui.messages.Get("success")
	context := ui.messages.Get("success_context")

	result := fmt.Sprintf("%s\n\n%s", title, message)
	if context != "" {
		result += fmt.Sprintf(" [%s]", context)
	}

	controls := "\nPress Enter to continue, q to quit"
	return result + controls
}

func (ui *UIController) errorView() string {
	title := errorStyle.Render("Error!")
	message := ui.messages.Get("error")
	context := ui.messages.Get("error_context")

	result := fmt.Sprintf("%s\n\n%s", title, message)
	if context != "" {
		result += fmt.Sprintf(" [%s]", context)
	}

	controls := "\nPress Enter to try again, q to quit"
	return result + controls
}

func (ui *UIController) handleResultInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		if ui.state == stateSuccess {
			if ui.isAuthenticated {
				ui.state = stateMenuLoggedIn
				ui.loggedInMenu = 0
				ui.messages.ClearAll()
				return ui, nil
			} else {
				ui.state = stateFinished
				return ui, tea.Quit
			}
		} else {
			if ui.isAuthenticated {
				ui.state = stateMenuLoggedIn
				ui.loggedInMenu = 0
			} else {
				ui.state = stateMenuLoggedOut
				ui.currentMenu = 0
				ui.login = ""
			}

			ui.messages.ClearAll()
			return ui, nil
		}
	}
	return ui, nil
}

func (ui *UIController) handleProcessComplete(msg processComplete) (tea.Model, tea.Cmd) {
	if msg.success {
		switch msg.context {
		case "auth_to_master":
			ui.state = stateMasterPassword
			ui.input = ""
			return ui, nil
		case "master_password":
			ui.state = stateMenuLoggedIn
			ui.input = ""
			return ui, nil
		case "auth":
			ui.state = stateMenuLoggedIn
			ui.input = ""
			return ui, nil
		case "delete_item":
			ui.state = stateDeleteSuccess
			ui.deleteSuccessMsg = msg.message
			ui.selectedItem = nil
			ui.decryptedItem = nil
			return ui, nil
		case "edit_item":
			ui.state = stateEditSuccess
			ui.editSuccessMsg = msg.message
			return ui, nil
		case "save_item":
			ui.state = stateAddItemSuccess
			ui.addItemSuccessMsg = msg.message
			ui.newItem = models.Item{}
			return ui, nil
		case "save_metadata":
			ui.state = stateMetadataSuccess
			ui.metadataSuccessMsg = msg.message
			return ui, nil
		case "logout":
			ui.state = stateLogoutSuccess
			ui.logoutSuccessMsg = msg.message
			return ui, nil
		default:
			ui.state = stateMenuLoggedIn
			ui.input = ""
			return ui, nil
		}
	} else {
		switch msg.context {
		case "auth":
			ui.state = stateMenuLoggedOut
		case "master_password":
			ui.state = stateMasterPassword
		case "delete_item":
			ui.state = stateDeleteError
			ui.deleteErrorMsg = msg.message
		case "edit_item":
			ui.state = stateEditError
			ui.editErrorMsg = msg.message
		case "save_item":
			ui.state = stateAddItemError
			ui.addItemErrorMsg = msg.message
		case "save_metadata":
			ui.state = stateMetadataError
			ui.metadataErrorMsg = msg.message
			return ui, nil
		case "logout":
			ui.state = stateLogoutError
			ui.logoutErrorMsg = msg.message
			return ui, nil
		default:
			ui.state = stateMenuLoggedOut
		}

		ui.input = ""
		return ui, nil
	}
}
