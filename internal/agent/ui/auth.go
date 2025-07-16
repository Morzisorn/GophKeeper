package ui

import (
	"context"
	"fmt"
	"gophkeeper/models"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (ui *UIController) handleMenuLoggedOutInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return ui, tea.Quit
	case "up", "k":
		if ui.currentMenu > 0 {
			ui.currentMenu--
		}
	case "down", "j":
		if ui.currentMenu < 1 {
			ui.currentMenu++
		}
	case "1":
		ui.currentMenu = 0
		return ui.startSignUp(), nil
	case "2":
		ui.currentMenu = 1
		return ui.startSignIn(), nil
	case "enter":
		if ui.currentMenu == 0 {
			return ui.startSignUp(), nil
		} else {
			return ui.startSignIn(), nil
		}
	}
	return ui, nil
}

func (ui *UIController) handleLoginInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return ui, tea.Quit
	case "esc":
		ui.state = stateMenuLoggedOut
		ui.input = ""
		ui.messages.ClearAll()
		return ui, nil
	case "enter":
		ui.login = strings.TrimSpace(ui.input)
		ui.input = ""

		// Переходим к вводу пароля в зависимости от текущего состояния
		if ui.state.IsSignUp() {
			ui.state = stateSignUpPassword
		} else {
			ui.state = stateSignInPassword
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

func (ui *UIController) handlePasswordInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return ui, tea.Quit
	case "esc":
		ui.input = ui.login

		// Return to input login depends on state
		if ui.state.IsSignUp() {
			ui.state = stateSignUpLogin
		} else {
			ui.state = stateSignInLogin
		}
		return ui, nil
	case "enter":
		password := strings.TrimSpace(ui.input)
		ui.input = ""
		ui.state = stateProcessing

		// Choose cmd depends on state
		if ui.state.IsSignUp() {
			return ui, ui.signUpCmd(ui.login, password)
		} else {
			return ui, ui.signInCmd(ui.login, password)
		}
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

func (ui *UIController) startSignUp() *UIController {
	ui.state = stateSignUpLogin
	ui.input = ""
	ui.messages.ClearAll()
	return ui
}

func (ui *UIController) startSignIn() *UIController {
	ui.state = stateSignInLogin
	ui.input = ""
	ui.messages.ClearAll()
	return ui
}

func (ui *UIController) signUpCmd(login, password string) tea.Cmd {
	return func() tea.Msg {
		err := ui.User.SignUpUser(context.Background(), &models.User{Login: login, Password: []byte(password)})
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Sign up error: %v", err),
				context: "auth",
			}
		}
		if ui.User.Client.GetJWTToken() != "" {
			ui.isAuthenticated = true
			ui.login = login
		}
		return processComplete{
			success: true,
			message: "Account successfully registered",
			context: "auth",
		}
	}
}

func (ui *UIController) signInCmd(login, password string) tea.Cmd {
	return func() tea.Msg {
		err := ui.User.SignInUser(context.Background(), &models.User{Login: login, Password: []byte(password)})
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Sign in error: %v", err),
				context: "auth",
			}
		}

		if ui.User.Client.GetJWTToken() != "" {
			ui.isAuthenticated = true
			ui.login = login
		}

		return processComplete{
			success: true,
			message: "Signed in successfully",
			context: "auth",
		}
	}
}

func (ui *UIController) loginInputView() string {
	var action string
	if ui.state.IsSignUp() {
		action = "Sign Up"
	} else {
		action = "Sign In"
	}

	title := titleStyle.Render(fmt.Sprintf("%s - Enter Login", action))
	input := inputStyle.Render(ui.input + "█")
	controls := "\nControls: Esc to go back, Enter to continue"

	return fmt.Sprintf("%s\n\nLogin: %s%s", title, input, controls)
}

func (ui *UIController) passwordInputView() string {
	var action string
	if ui.state.IsSignUp() {
		action = "Sign Up"
	} else {
		action = "Sign In"
	}

	title := titleStyle.Render(fmt.Sprintf("%s - Enter Password", action))

	// Hide password with ***
	hiddenPassword := strings.Repeat("*", len(ui.input))
	input := inputStyle.Render(hiddenPassword + "█")
	controls := "\nControls: Esc to go back, Enter to continue"

	return fmt.Sprintf("%s\n\nLogin: %s\nPassword: %s%s", title, ui.login, input, controls)
}

func (ui *UIController) processingView() string {
	var action string
	if ui.state.IsSignUp() {
		action = "Signing up..."
	} else {
		action = "Signing in..."
	}

	return fmt.Sprintf("\n%s\n\nPlease wait...", action)
}


func (ui *UIController) handleMasterPasswordInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return ui, tea.Quit
	case "esc":
		ui.state = stateMenuLoggedIn
		ui.input = ""
		ui.messages.ClearAll()
		return ui, nil
	case "enter":
		masterPassword := strings.TrimSpace(ui.input)
		ui.input = ""
		ui.state = stateProcessing

		return ui, ui.setMasterPasswordCmd(masterPassword)
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

func (ui *UIController) setMasterPasswordCmd(masterPassword string) tea.Cmd {
	return func() tea.Msg {
		ui.User.SetMasterKey(masterPassword)

		return processComplete{
			success: true,
			message: "Master password set. Welcome!",
			context: "master_password",
		}
	}
}