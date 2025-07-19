package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_handleMenuLoggedOutInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Тест ctrl+c
	model, cmd := ui.handleMenuLoggedOutInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Тест q
	model, cmd = ui.handleMenuLoggedOutInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleMenuLoggedOutInput_Navigation(t *testing.T) {
	ui := &UIController{
		menuCtrl: menuCtrl{
			currentMenu: 1,
		},
	}

	// Тест up
	model, cmd := ui.handleMenuLoggedOutInput(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.currentMenu)

	// Тест down
	model, cmd = ui.handleMenuLoggedOutInput(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.currentMenu)
}

func TestUIController_handleMenuLoggedOutInput_DirectSelection(t *testing.T) {
	ui := &UIController{}

	// Тест выбора 1 (Sign Up)
	model, cmd := ui.handleMenuLoggedOutInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 0, ui.currentMenu)

	// Тест выбора 2 (Sign In)
	model, cmd = ui.handleMenuLoggedOutInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, 1, ui.currentMenu)
}

func TestUIController_handleLoginInput_Quit(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleLoginInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleLoginInput_Escape(t *testing.T) {
	ui := &UIController{
		input: "test-input",
		state: stateSignUpLogin,
	}

	model, cmd := ui.handleLoginInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMenuLoggedOut, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_handleLoginInput_Enter(t *testing.T) {
	ui := &UIController{
		input: "  test-login  ",
		state: stateSignUpLogin,
	}

	model, cmd := ui.handleLoginInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "test-login", ui.login)
	assert.Empty(t, ui.input)
	assert.Equal(t, stateSignUpPassword, ui.state)
}

func TestUIController_handleLoginInput_Backspace(t *testing.T) {
	ui := &UIController{input: "test"}

	model, cmd := ui.handleLoginInput(tea.KeyMsg{Type: tea.KeyBackspace})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "tes", ui.input)
}

func TestUIController_handleLoginInput_CharacterInput(t *testing.T) {
	ui := &UIController{input: "test"}

	model, cmd := ui.handleLoginInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "testa", ui.input)
}

func TestUIController_handlePasswordInput_Quit(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handlePasswordInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handlePasswordInput_Escape(t *testing.T) {
	ui := &UIController{
		userCtrl: userCtrl{
			login: "test-login",
		},
		input: "password",
		state: stateSignUpPassword,
	}

	model, cmd := ui.handlePasswordInput(tea.KeyMsg{Type: tea.KeyEscape})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "test-login", ui.input)
	assert.Equal(t, stateSignUpLogin, ui.state)
}

func TestUIController_handlePasswordInput_Enter_SignUp(t *testing.T) {
	ui := &UIController{
		input: "  password  ",
		state: stateSignUpPassword,
		userCtrl: userCtrl{
			login: "test-login",
		},
	}

	model, cmd := ui.handlePasswordInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
	assert.Empty(t, ui.input)
	assert.Equal(t, stateProcessing, ui.state)
}

func TestUIController_handlePasswordInput_Enter_SignIn(t *testing.T) {
	ui := &UIController{
		input: "  password  ",
		state: stateSignInPassword,
		userCtrl: userCtrl{
			login: "test-login",
		},
	}

	model, cmd := ui.handlePasswordInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
	assert.Empty(t, ui.input)
	assert.Equal(t, stateProcessing, ui.state)
}

func TestUIController_startSignUp(t *testing.T) {
	ui := &UIController{
		input: "test-input",
		state: stateMenuLoggedOut,
	}

	result := ui.startSignUp()

	assert.Equal(t, ui, result)
	assert.Equal(t, stateSignUpLogin, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_startSignIn(t *testing.T) {
	ui := &UIController{
		input: "test-input",
		state: stateMenuLoggedOut,
	}

	result := ui.startSignIn()

	assert.Equal(t, ui, result)
	assert.Equal(t, stateSignInLogin, ui.state)
	assert.Empty(t, ui.input)
}

func TestUIController_handleMasterPasswordInput_Quit(t *testing.T) {
	ui := &UIController{}

	model, cmd := ui.handleMasterPasswordInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

// func TestUIController_handleMasterPasswordInput_Escape(t *testing.T) {
// 	ui := &UIController{
// 		input: "test-input",
// 		state: stateMasterPasswordInput,
// 	}

// 	model, cmd := ui.handleMasterPasswordInput(tea.KeyMsg{Type: tea.KeyEscape})

// 	assert.Equal(t, ui, model)
// 	assert.Nil(t, cmd)
// 	assert.Equal(t, stateMenuLoggedIn, ui.state)
// 	assert.Empty(t, ui.input)
// }

// func TestUIController_handleMasterPasswordInput_Enter(t *testing.T) {
// 	ui := &UIController{
// 		input: "  master-password  ",
// 		state: stateMasterPasswordInput,
// 	}

// 	model, cmd := ui.handleMasterPasswordInput(tea.KeyMsg{Type: tea.KeyEnter})

// 	assert.Equal(t, ui, model)
// 	assert.NotNil(t, cmd)
// 	assert.Empty(t, ui.input)
// 	assert.Equal(t, stateProcessing, ui.state)
// }

func TestUIController_handleMasterPasswordInput_Backspace(t *testing.T) {
	ui := &UIController{input: "test"}

	model, cmd := ui.handleMasterPasswordInput(tea.KeyMsg{Type: tea.KeyBackspace})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "tes", ui.input)
}

func TestUIController_handleMasterPasswordInput_CharacterInput(t *testing.T) {
	ui := &UIController{input: "test"}

	model, cmd := ui.handleMasterPasswordInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, "testa", ui.input)
}

// Тесты view методов
func TestUIController_loginInputView_SignUp(t *testing.T) {
	ui := &UIController{
		state: stateSignUpLogin,
		input: "test-login",
	}

	view := ui.loginInputView()

	assert.Contains(t, view, "Sign Up")
	assert.Contains(t, view, "test-login")
	assert.Contains(t, view, "█")
}

func TestUIController_loginInputView_SignIn(t *testing.T) {
	ui := &UIController{
		state: stateSignInLogin,
		input: "test-login",
	}

	view := ui.loginInputView()

	assert.Contains(t, view, "Sign In")
	assert.Contains(t, view, "test-login")
	assert.Contains(t, view, "█")
}

func TestUIController_passwordInputView_SignUp(t *testing.T) {
	ui := &UIController{
		state: stateSignUpPassword,
		input: "password",
		userCtrl: userCtrl{
			login: "test-login",
		},
	}

	view := ui.passwordInputView()

	assert.Contains(t, view, "Sign Up")
	assert.Contains(t, view, "test-login")
	assert.Contains(t, view, "********")
	assert.NotContains(t, view, "password")
}

func TestUIController_passwordInputView_SignIn(t *testing.T) {
	ui := &UIController{
		state: stateSignInPassword,
		input: "password",
		userCtrl: userCtrl{
			login: "test-login",
		},
	}

	view := ui.passwordInputView()

	assert.Contains(t, view, "Sign In")
	assert.Contains(t, view, "test-login")
	assert.Contains(t, view, "********")
	assert.NotContains(t, view, "password")
}

func TestUIController_masterPasswordInputView(t *testing.T) {
	ui := &UIController{
		input: "masterpass",
	}

	view := ui.masterPasswordInputView()

	assert.Contains(t, view, "Master Password Required")
	assert.Contains(t, view, "**********")
	assert.NotContains(t, view, "masterpass")
	assert.Contains(t, view, "█")
}

func TestUIController_processingView(t *testing.T) {
	ui := &UIController{}

	view := ui.processingView()

	assert.Contains(t, view, "Processing")
	assert.Contains(t, view, "Please wait")
}
