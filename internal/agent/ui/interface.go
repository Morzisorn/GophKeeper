package ui

import (
	"bufio"
	"context"
	"fmt"
	"gophkeeper/internal/agent/services"
	"os"
	"strings"
)

type UserInterface interface {
	Run() error
}

type UIController struct {
	scanner *bufio.Scanner
	User    *services.UserService
}

const (
	messageWelcome      = "Welcome to GophKeeper - best sensetive info manager :)"
	messageChooseOption = "Choose an option - enter number or command:"
	optionsAnauthorized = `1. Sign up
2. Sign in`
	optionInvalid        = "Invalid option :("
	messageEnterLogin    = "Enter login:"
	messageEnterPassword = "Enter password:"
)

func NewUIController(us *services.UserService) UserInterface {
	return &UIController{
		User:    us,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (ui *UIController) Run() error {
	for {
		if ui.User.Client.GetJWTToken() == "" {
			err := ui.anauthorizedFlow()
			if err != nil {
				return err
			}

			fmt.Println("Ready for menu")
			break
		}

	}
	fmt.Println("Finish")
	return nil
}

func (ui *UIController) anauthorizedFlow() error {
	fmt.Println(messageWelcome)
	fmt.Println(messageChooseOption)
	fmt.Println(optionsAnauthorized)
	if !ui.scanner.Scan() {
		return nil
	}

	message := strings.ToLower(strings.TrimSpace(ui.scanner.Text()))
	switch {
	case message == "1" || message == "sign up":
		if err := ui.signUpUser(); err != nil {
			return err
		}
		fmt.Println("Account successfully registered")
	case message == "2" || message == "sign in":
		if err := ui.signInUser(); err != nil {
			return err
		}
		fmt.Println("Signed in successfully")
	}
	return nil
}

func (ui *UIController) getLoginPassword() (login, password string) {
	fmt.Println(messageEnterLogin)
	ui.scanner.Scan()
	login = strings.TrimSpace(ui.scanner.Text())
	fmt.Println(messageEnterPassword)
	ui.scanner.Scan()
	password = strings.TrimSpace(ui.scanner.Text())
	return login, password
}

func (ui *UIController) signUpUser() error {
	login, password := ui.getLoginPassword()
	err := ui.User.SignUpUser(context.Background(), login, password)
	if err != nil {
		return fmt.Errorf("sign up user error: %w", err)
	}
	return nil
}

func (ui *UIController) signInUser() error {
	login, password := ui.getLoginPassword()
	err := ui.User.SignInUser(context.Background(), login, password)
	if err != nil {
		return fmt.Errorf("sign in user error: %w", err)
	}
	return nil
}
