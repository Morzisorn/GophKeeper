package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUIController_menuLoggedOutView(t *testing.T) {
	ui := &UIController{
		menuCtrl: menuCtrl{
			currentMenu: 0,
		},
	}
	
	view := ui.menuLoggedOutView()
	
	assert.Contains(t, view, "Welcome to GophKeeper")
	assert.Contains(t, view, "Sign up")
	assert.Contains(t, view, "Sign in")
	assert.Contains(t, view, "1.")
	assert.Contains(t, view, "2.")
	assert.NotEmpty(t, view)
}