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
	User              *services.UserService
	Item              *services.ItemService
	state             state
	currentMenu       int
	input             string
	login             string
	messages          Messages
	isAuthenticated   bool
	loggedInMenu      int
	maxLoggedInMenu   int
	items             []models.EncryptedItem
	currentItem       int
	selectedItem      *models.EncryptedItem
	decryptedItem     *models.Item
	maxItems          int
	newItem           models.Item
	addItemErrorMsg   string 
	addItemSuccessMsg string
	itemTypeMenu      int 
	maxItemTypes      int 
	confirmChoice     int
	decryptErrorMsg   string
	deleteSuccessMsg  string
	deleteErrorMsg    string
	editingItem       *models.Item 
	editStep          int          
	editSuccessMsg    string       
	editErrorMsg      string
	metadataMenu       int  
	metadataKeys       []string 
	currentMetaKey     string   
	currentMetaValue   string   
	metadataSuccessMsg string   
	metadataErrorMsg   string

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
