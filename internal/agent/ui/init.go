package ui

import (
	"fmt"
	"gophkeeper/internal/agent/services"
	"gophkeeper/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UserInterface interface {
	Run() error
}

type UIController struct {
	User            *services.UserService
	Item            *services.ItemService
	state           state
	currentMenu     int
	input           string
	login           string
	messages        Messages
	isAuthenticated bool
	loggedInMenu    int
	maxLoggedInMenu int
	items           []models.Item
	currentItem     int
	maxItems        int
	newItem         models.Item
	itemTypeMenu    int // Выбранный тип элемента
	maxItemTypes    int // Максимальный индекс типов
}

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

func NewUIController(us *services.UserService, is *services.ItemService) UserInterface {
	ui := &UIController{
		User:            us,
		Item:            is,
		state:           stateMenuLoggedOut,
		maxLoggedInMenu: 4,
	}
	ui.messages.init()
	return ui
}

func (ui *UIController) Run() error {
	p := tea.NewProgram(ui, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	controller := finalModel.(*UIController)
	if controller.state == stateError {
		return fmt.Errorf(controller.messages.Get("error"))
	}

	return nil
}

func (ui *UIController) Init() tea.Cmd {
	return nil
}

func (ui *UIController) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return ui.handleKeyMsg(msg)
	case itemsLoaded:
		ui.items = msg.items
		ui.maxItems = len(msg.items) - 1
		if len(msg.items) == 0 {
			ui.messages.Set("info", "No items found")
		} else {
			ui.messages.Clear("info")
		}
		return ui, nil
	case processComplete:
		if msg.success {
			ui.state = stateSuccess
			ui.messages.Set("success", msg.message)
			if msg.context != "" {
				ui.messages.Set("success_context", msg.context)
			}
		} else {
			ui.state = stateError
			ui.messages.Set("error", msg.message)
			if msg.context != "" {
				ui.messages.Set("error_context", msg.context)
			}
		}
		return ui, nil
	case errorMsg:
		ui.state = stateError
		ui.messages.Set("error", msg.err.Error())
		if msg.context != "" {
			ui.messages.Set("error_context", msg.context)
		}
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
	case ui.state.IsAddItemInput():
		return ui.handleAddItemInput(msg)
	case ui.state.IsLoginInput():
		return ui.handleLoginInput(msg)
	case ui.state.IsPasswordInput():
		return ui.handlePasswordInput(msg)
	case ui.state == stateSuccess || ui.state == stateError:
		return ui.handleResultInput(msg)
	case ui.state == stateProcessing:
		// Ignore input during processing
		return ui, nil
	}
	return ui, nil
}

func (ui *UIController) handleResultInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "enter":
		if ui.state == stateSuccess {
			// При успехе проверяем авторизацию
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
			// При ошибке возвращаемся в соответствующее меню
			//context := ui.messages.Get("error_context")

			if ui.isAuthenticated {
				// Если пользователь авторизован, возвращаемся в меню авторизованного
				ui.state = stateMenuLoggedIn
				ui.loggedInMenu = 0
			} else {
				// Если не авторизован, возвращаемся в меню неавторизованного
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

func (ui *UIController) View() string {
	debug := fmt.Sprintf("Current state: %d\n", int(ui.state))
	switch {
	case ui.state == stateMenuLoggedOut:
		return ui.menuLoggedOutView()
	case ui.state == stateMenuLoggedIn:
		return ui.menuLoggedInView()
	case ui.state == stateItemsList:
		return ui.itemsListView()
	case ui.state.IsLoginInput():
		return ui.loginInputView()
	case ui.state.IsPasswordInput():
		return ui.passwordInputView()
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
	}
	return "View error:" + debug
}

func (ui *UIController) menuLoggedOutView() string {
	title := titleStyle.Render("Welcome to GophKeeper - best sensitive info manager :)")
	subtitle := "Choose an option - enter number or use arrow keys:"

	options := []string{
		"Sign up",
		"Sign in",
	}

	menu := ""
	for i, option := range options {
		prefix := fmt.Sprintf("%d. ", i+1)
		if i == ui.currentMenu {
			menu += selectedStyle.Render(prefix+option) + "\n"
		} else {
			menu += menuStyle.Render(prefix+option) + "\n"
		}
	}

	controls := "\nControls: ↑/↓ to navigate, Enter to select, q to quit"

	return fmt.Sprintf("%s\n\n%s\n\n%s%s", title, subtitle, menu, controls)
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
