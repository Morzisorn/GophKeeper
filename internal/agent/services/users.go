package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"gophkeeper/internal/agent/client"
	"gophkeeper/internal/errs"
)

type UserService struct {
	Client client.Client
}

func NewUserService(client client.Client) *UserService {
	return &UserService{
		Client: client,
	}
}

func (us *UserService) SignUpUser(ctx context.Context, login, password string) error {
	if login == "" || password == "" {
		return errs.ErrRequiredArgumentIsMissing
	}

	encryptedPassword, err := encryptData([]byte(password))
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	token, err := us.Client.SignUpUser(ctx, login, base64.StdEncoding.EncodeToString(encryptedPassword))
	if err != nil {
		return fmt.Errorf("server failed to sign up user: %w", err)
	}

	us.Client.SetJWTToken(token)

	return nil
}

func (us *UserService) SignInUser(ctx context.Context, login, password string) error {
	if login == "" || password == "" {
		return errs.ErrRequiredArgumentIsMissing
	}

	encryptedPassword, err := encryptData([]byte(password))
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	token, err := us.Client.SignInUser(ctx, login, base64.StdEncoding.EncodeToString(encryptedPassword))
	if err != nil {
		return fmt.Errorf("server failed to sign in user: %w", err)
	}

	us.Client.SetJWTToken(token)

	return nil
}