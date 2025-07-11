package ui

import (
	"context"
	"fmt"
	"gophkeeper/internal/agent/services"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type UserInterface interface {
	Run() error
}

type UIController struct {
	User *services.UserService
}

// App state
type state int

const (
	stateMenu state = iota
	stateSignUp
	stateSignIn
	stateLoginInput
	statePasswordInput
	stateProcessing
	stateSuccess
	stateError
	stateFinished
)

// Модель приложения
type model struct {
	state        state
	currentMenu  int
	input        string
	login        string
	password     string
	errorMsg     string
	successMsg   string
	userService  *services.UserService
	isSignUp     bool
	cursorPos    int
	showPassword bool
}

// Messages
type (
	processComplete struct {
		success bool
		message string
	}
	errorMsg struct {
		err error
	}
)

// Styles
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

func NewUIController(us *services.UserService) UserInterface {
	return &UIController{
		User: us,
	}
}

func (ui *UIController) Run() error {
	initialModel := model{
		state:       stateMenu,
		userService: ui.User,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	// Check final state
	m := finalModel.(model)
	if m.state == stateError {
		return fmt.Errorf(m.errorMsg)
	}

	return nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case processComplete:
		if msg.success {
			m.state = stateSuccess
			m.successMsg = msg.message
		} else {
			m.state = stateError
			m.errorMsg = msg.message
		}
		return m, nil
	case errorMsg:
		m.state = stateError
		m.errorMsg = msg.err.Error()
		return m, nil
	}
	return m, nil
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateMenu:
		return m.handleMenuInput(msg)
	case stateLoginInput:
		return m.handleLoginInput(msg)
	case statePasswordInput:
		return m.handlePasswordInput(msg)
	case stateSuccess, stateError:
		return m.handleResultInput(msg)
	case stateProcessing:
		// Ignore input during processing
		return m, nil
	}
	return m, nil
}

func (m model) handleMenuInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.currentMenu > 0 {
			m.currentMenu--
		}
	case "down", "j":
		if m.currentMenu < 1 {
			m.currentMenu++
		}
	case "1":
		m.currentMenu = 0
		return m.startSignUp(), nil
	case "2":
		m.currentMenu = 1
		return m.startSignIn(), nil
	case "enter":
		if m.currentMenu == 0 {
			return m.startSignUp(), nil
		} else {
			return m.startSignIn(), nil
		}
	}
	return m, nil
}

func (m model) handleLoginInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = stateMenu
		m.input = ""
		return m, nil
	case "enter":
		m.login = strings.TrimSpace(m.input)
		m.input = ""
		m.state = statePasswordInput
		return m, nil
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) handlePasswordInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.state = stateLoginInput
		m.input = m.login
		return m, nil
	case "enter":
		m.password = strings.TrimSpace(m.input)
		m.input = ""
		m.state = stateProcessing
		
		if m.isSignUp {
			return m, m.signUpCmd()
		} else {
			return m, m.signInCmd()
		}
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m model) handleResultInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "enter":
		if m.state == stateSuccess {
			m.state = stateFinished
			return m, tea.Quit
		} else {
			// При ошибке возвращаемся в меню
			m.state = stateMenu
			m.errorMsg = ""
			m.login = ""
			m.password = ""
			return m, nil
		}
	}
	return m, nil
}

func (m model) startSignUp() model {
	m.isSignUp = true
	m.state = stateLoginInput
	m.input = ""
	return m
}

func (m model) startSignIn() model {
	m.isSignUp = false
	m.state = stateLoginInput
	m.input = ""
	return m
}

func (m model) signUpCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.userService.SignUpUser(context.Background(), m.login, m.password)
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Sign up error: %v", err),
			}
		}
		return processComplete{
			success: true,
			message: "Account successfully registered",
		}
	}
}

func (m model) signInCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.userService.SignInUser(context.Background(), m.login, m.password)
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Sign in error: %v", err),
			}
		}
		return processComplete{
			success: true,
			message: "Signed in successfully",
		}
	}
}

func (m model) View() string {
	switch m.state {
	case stateMenu:
		return m.menuView()
	case stateLoginInput:
		return m.loginInputView()
	case statePasswordInput:
		return m.passwordInputView()
	case stateProcessing:
		return m.processingView()
	case stateSuccess:
		return m.successView()
	case stateError:
		return m.errorView()
	}
	return ""
}

func (m model) menuView() string {
	title := titleStyle.Render("Welcome to GophKeeper - best sensitive info manager :)")
	subtitle := "Choose an option - enter number or use arrow keys:"
	
	options := []string{
		"Sign up",
		"Sign in",
	}
	
	menu := ""
	for i, option := range options {
		prefix := fmt.Sprintf("%d. ", i+1)
		if i == m.currentMenu {
			menu += selectedStyle.Render(prefix+option) + "\n"
		} else {
			menu += menuStyle.Render(prefix+option) + "\n"
		}
	}
	
	controls := "\nControls: ↑/↓ to navigate, Enter to select, q to quit"
	
	return fmt.Sprintf("%s\n\n%s\n\n%s%s", title, subtitle, menu, controls)
}

func (m model) loginInputView() string {
	action := "Sign In"
	if m.isSignUp {
		action = "Sign Up"
	}
	
	title := titleStyle.Render(fmt.Sprintf("%s - Enter Login", action))
	input := inputStyle.Render(m.input + "█")
	controls := "\nControls: Esc to go back, Enter to continue"
	
	return fmt.Sprintf("%s\n\nLogin: %s%s", title, input, controls)
}

func (m model) passwordInputView() string {
	action := "Sign In"
	if m.isSignUp {
		action = "Sign Up"
	}
	
	title := titleStyle.Render(fmt.Sprintf("%s - Enter Password", action))
	
	// Скрываем пароль звездочками
	hiddenPassword := strings.Repeat("*", len(m.input))
	input := inputStyle.Render(hiddenPassword + "█")
	controls := "\nControls: Esc to go back, Enter to continue"
	
	return fmt.Sprintf("%s\n\nLogin: %s\nPassword: %s%s", title, m.login, input, controls)
}

func (m model) processingView() string {
	action := "Signing in..."
	if m.isSignUp {
		action = "Signing up..."
	}
	
	return fmt.Sprintf("\n%s\n\nPlease wait...", action)
}

func (m model) successView() string {
	title := successStyle.Render("Success!")
	message := m.successMsg
	controls := "\nPress Enter to continue, q to quit"
	
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}

func (m model) errorView() string {
	title := errorStyle.Render("Error!")
	message := m.errorMsg
	controls := "\nPress Enter to try again, q to quit"
	
	return fmt.Sprintf("%s\n\n%s%s", title, message, controls)
}