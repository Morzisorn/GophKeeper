package ui

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestUIController_decryptItemCmd_Success(t *testing.T) {
	// This test is complex for unit testing as it requires mock Item service
	// Skipping as this is an integration test
	t.Skip("Requires mock Item service - integration test")
}

func TestUIController_handleDecryptError_IncorrectPassword(t *testing.T) {
	ui := &UIController{}

	msg := decryptError{
		err:     errors.New("incorrect master password"),
		context: "decrypt_item",
	}

	model, cmd := ui.handleDecryptError(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateDecryptError, ui.state)
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "Incorrect master password")
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "restart the application")
}

func TestUIController_handleDecryptError_AuthenticationFailed(t *testing.T) {
	ui := &UIController{}

	msg := decryptError{
		err:     errors.New("cipher: message authentication failed"),
		context: "decrypt_item",
	}

	model, cmd := ui.handleDecryptError(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateDecryptError, ui.state)
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "Failed to decrypt item")
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "cipher: message authentication failed")
}

func TestUIController_handleDecryptError_FailedToDecrypt(t *testing.T) {
	ui := &UIController{}

	msg := decryptError{
		err:     errors.New("failed to decrypt data"),
		context: "decrypt_item",
	}

	model, cmd := ui.handleDecryptError(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateDecryptError, ui.state)
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "Failed to decrypt item")
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "failed to decrypt data")
}

func TestUIController_handleDecryptError_DecryptItemCmdGeneratedError(t *testing.T) {
	// This test checks the case when decryptItemCmd generates "incorrect master password" error
	ui := &UIController{}

	msg := decryptError{
		err:     errors.New("incorrect master password"), // As created in decryptItemCmd
		context: "decrypt_item",
	}

	model, cmd := ui.handleDecryptError(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateDecryptError, ui.state)
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "Incorrect master password")
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "restart the application")
}

func TestUIController_handleDecryptError_OtherError(t *testing.T) {
	ui := &UIController{}

	msg := decryptError{
		err:     errors.New("some other error"),
		context: "decrypt_item",
	}

	model, cmd := ui.handleDecryptError(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateDecryptError, ui.state)
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "Failed to decrypt item")
	assert.Contains(t, ui.itemCtrl.decryptErrorMsg, "some other error")
}

func TestUIController_handleDecryptError_DifferentContext(t *testing.T) {
	ui := &UIController{}

	msg := decryptError{
		err:     errors.New("some error"),
		context: "other_context",
	}

	model, cmd := ui.handleDecryptError(msg)

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateItemsList, ui.state)
}

func TestUIController_handleDecryptErrorInput_Quit(t *testing.T) {
	ui := &UIController{}

	// Test ctrl+c
	model, cmd := ui.handleDecryptErrorInput(tea.KeyMsg{Type: tea.KeyCtrlC})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)

	// Test q
	model, cmd = ui.handleDecryptErrorInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	assert.Equal(t, ui, model)
	assert.NotNil(t, cmd)
}

func TestUIController_handleDecryptErrorInput_Enter(t *testing.T) {
	ui := &UIController{
		input: "test-input",
		itemCtrl: itemCtrl{
			decryptErrorMsg: "test error message",
		},
	}

	model, cmd := ui.handleDecryptErrorInput(tea.KeyMsg{Type: tea.KeyEnter})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateMasterPassword, ui.state)
	assert.Empty(t, ui.input)
	assert.Empty(t, ui.itemCtrl.decryptErrorMsg)
}

func TestUIController_handleDecryptErrorInput_OtherKey(t *testing.T) {
	ui := &UIController{
		state: stateDecryptError,
	}

	model, cmd := ui.handleDecryptErrorInput(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})

	assert.Equal(t, ui, model)
	assert.Nil(t, cmd)
	assert.Equal(t, stateDecryptError, ui.state) // Should remain unchanged
}

func TestUIController_decryptErrorView_IncorrectPassword(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptErrorMsg: "Incorrect master password. Please restart the application and enter the correct master password.",
		},
	}

	view := ui.decryptErrorView()

	assert.Contains(t, view, "Decryption Error")
	assert.Contains(t, view, "Incorrect master password")
	assert.Contains(t, view, "restart the application")
	assert.Contains(t, view, "Enter to re-enter master password")
	assert.Contains(t, view, "q to quit")
}

func TestUIController_decryptErrorView_OtherError(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptErrorMsg: "Failed to decrypt item: some other error",
		},
	}

	view := ui.decryptErrorView()

	assert.Contains(t, view, "Decryption Error")
	assert.Contains(t, view, "Failed to decrypt item")
	assert.Contains(t, view, "some other error")
	assert.Contains(t, view, "Enter to re-enter master password")
	assert.Contains(t, view, "q to quit")
}

func TestUIController_decryptErrorView_EmptyMessage(t *testing.T) {
	ui := &UIController{
		itemCtrl: itemCtrl{
			decryptErrorMsg: "",
		},
	}

	view := ui.decryptErrorView()

	assert.Contains(t, view, "Decryption Error")
	assert.Contains(t, view, "Enter to re-enter master password")
	assert.Contains(t, view, "q to quit")
}

func TestUIController_decryptErrorView_ControlsAreSame(t *testing.T) {
	// Check that controls are the same for both error types
	ui1 := &UIController{
		itemCtrl: itemCtrl{
			decryptErrorMsg: "incorrect master password detected",
		},
	}

	ui2 := &UIController{
		itemCtrl: itemCtrl{
			decryptErrorMsg: "some other decrypt error",
		},
	}

	view1 := ui1.decryptErrorView()
	view2 := ui2.decryptErrorView()

	// Both should contain the same controls
	assert.Contains(t, view1, "Enter to re-enter master password")
	assert.Contains(t, view2, "Enter to re-enter master password")
	assert.Contains(t, view1, "q to quit")
	assert.Contains(t, view2, "q to quit")
}
