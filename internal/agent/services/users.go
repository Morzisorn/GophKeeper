package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"gophkeeper/config"
	"gophkeeper/internal/agent/client"
	"gophkeeper/internal/errs"
	"gophkeeper/models"
)

type UserService struct {
	Client client.Client
	Crypto *CryptoService
	config *config.Config
}

func NewUserService(client client.Client, cr *CryptoService) *UserService {
	return &UserService{
		Client: client,
		Crypto: cr,
		config: config.GetAgentConfig(),
	}
}

func (us *UserService) SignUpUser(ctx context.Context, user *models.User) error {
	if user.Login == "" || user.Password == nil {
		return errs.ErrRequiredArgumentIsMissing
	}

	encryptedPassword, err := encryptData(user.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	base64.StdEncoding.Encode(user.Password, encryptedPassword)

	token, salt, err := us.Client.SignUpUser(ctx, user)
	if err != nil {
		return fmt.Errorf("server failed to sign up user: %w", err)
	}

	err = us.Crypto.SetSalt(salt)
	if err != nil {
		return err
	}

	us.Client.SetJWTToken(token)

	return nil
}

func (us *UserService) SignInUser(ctx context.Context, user *models.User) error {
	if user.Login == "" || user.Password == nil {
		return errs.ErrRequiredArgumentIsMissing
	}

	encryptedPassword, err := encryptData(user.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	base64.StdEncoding.Encode(user.Password, encryptedPassword)

	token, salt, err := us.Client.SignInUser(ctx, user)
	if err != nil {
		return fmt.Errorf("server failed to sign in user: %w", err)
	}
	err = us.Crypto.SetSalt(salt)
	if err != nil {
		return err
	}

	us.Client.SetJWTToken(token)

	return nil
}


func (us *UserService) SetMasterKey(masterPassword string) {
	masterKey := us.Crypto.GenerateMasterKey([]byte(masterPassword))
	us.Crypto.SetMasterKey(masterKey)
}