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

func NewUserService(client client.Client, cr *CryptoService) (*UserService, error) {
	cnfg, err := config.GetAgentConfig()
	if err != nil {
		return nil, fmt.Errorf("get agent config error: %w", err)
	}
	return &UserService{
		Client: client,
		Crypto: cr,
		config: cnfg,
	}, nil
}

func (us *UserService) SignUpUser(ctx context.Context, user *models.User) error {
	if user.Login == "" || user.Password == nil {
		return errs.ErrRequiredArgumentIsMissing
	}

	encryptedPassword, err := encryptData(user.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	user.Password = []byte(base64.StdEncoding.EncodeToString(encryptedPassword))

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

	user.Password = []byte(base64.StdEncoding.EncodeToString(encryptedPassword))

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
	us.Crypto.SetMasterPassword(masterPassword)
	masterKey := us.Crypto.GenerateMasterKey()
	us.Crypto.SetMasterKey(masterKey)
}

func (us *UserService) Logout() error {
	us.config.MasterKey = nil
	us.config.MasterPassword = ""
	us.config.Salt = nil
	return nil
}
