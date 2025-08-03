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

		if ui.state.IsSignUp() {
			ui.state = stateSignUpLogin
		} else {
			ui.state = stateSignInLogin
		}
		return ui, nil
	case "enter":
		password := strings.TrimSpace(ui.input)
		ui.input = ""
		isSignUp := ui.state.IsSignUp()
		ui.state = stateProcessing

		if isSignUp {
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
		token, err := ui.User.Client.GetJWTToken()
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Sign up error: %v", err),
				context: "auth",
			}
		}
		if len(token) != 0 {
			ui.isAuthenticated = true
			ui.login = login
		}
		return processComplete{
			success: true,
			message: "Account successfully registered",
			context: "auth_to_master",
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

		token, err := ui.User.Client.GetJWTToken()
		if err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Sign up error: %v", err),
				context: "auth",
			}
		}

		if len(token) != 0 {
			ui.isAuthenticated = true
			ui.login = login
		}

		return processComplete{
			success: true,
			message: "Signed in successfully",
			context: "auth_to_master",
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

	hiddenPassword := strings.Repeat("*", len(ui.input))
	input := inputStyle.Render(hiddenPassword + "█")
	controls := "\nControls: Esc to go back, Enter to continue"

	return fmt.Sprintf("%s\n\nLogin: %s\nPassword: %s%s", title, ui.login, input, controls)
}

func (ui *UIController) processingView() string {
	return fmt.Sprintln("\nProcessing...\n\nPlease wait...")
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
		if err := ui.User.SetMasterKey(masterPassword); err != nil {
			return processComplete{
				success: false,
				message: fmt.Sprintf("Set master password: %v", err),
				context: "master_password",
			}
		}

		return processComplete{
			success: true,
			message: "Master password set. Welcome!",
			context: "master_password",
		}
	}
}

func (ui *UIController) masterPasswordInputView() string {
	title := titleStyle.Render("Master Password Required")

	hiddenPassword := strings.Repeat("*", len(ui.input))
	input := inputStyle.Render(hiddenPassword + "█")
	controls := "\nControls: Esc to go back, Enter to continue"

	info := "\nEnter your master password to unlock your vault:"

	return fmt.Sprintf("%s%s\n\nMaster Password: %s%s", title, info, input, controls)
}
