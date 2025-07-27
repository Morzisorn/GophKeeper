package ui

import (
	"fmt"
	"gophkeeper/internal/agent/services"
	"gophkeeper/models"

	tea "github.com/charmbracelet/bubbletea"
)

type UserInterface interface {
	Run() error
}

type UIController struct {
	User *services.UserService
	Item *services.ItemService

	state           state
	input           string
	messages        Messages
	loggedInMenu    int
	maxLoggedInMenu int

	confirmChoice int

	menuCtrl
	userCtrl
	itemCtrl
	logoutCtrl
}

type menuCtrl struct {
	currentMenu int
}

type userCtrl struct {
	login           string
	isAuthenticated bool
}

type itemCtrl struct {
	items    []models.EncryptedItem
	maxItems int

	currentItem     int
	selectedItem    *models.EncryptedItem
	decryptedItem   *models.Item
	decryptErrorMsg string

	itemTypeMenu int
	maxItemTypes int

	newItem           models.Item
	addItemErrorMsg   string
	addItemSuccessMsg string

	editingItem    *models.Item
	editStep       int
	editSuccessMsg string
	editErrorMsg   string

	deleteSuccessMsg string
	deleteErrorMsg   string

	itemMetaCtrl
}

type itemMetaCtrl struct {
	metadataMenu       int
	metadataKeys       []string
	currentMetaKey     string
	currentMetaValue   string
	metadataSuccessMsg string
	metadataErrorMsg   string
}

type logoutCtrl struct {
	logoutSuccessMsg string
	logoutErrorMsg   string
}

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
