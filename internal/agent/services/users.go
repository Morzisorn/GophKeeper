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
	crypto *CryptoService
	cnfg   config.AgentUserServiceConfig
}

func NewUserService(cnfg config.AgentUserServiceConfig, client client.Client, cr *CryptoService) (*UserService, error) {
	return &UserService{
		Client: client,
		crypto: cr,
		cnfg:   cnfg,
	}, nil
}

func (us *UserService) SignUpUser(ctx context.Context, user *models.User) error {
	if user.Login == "" || user.Password == nil {
		return errs.ErrRequiredArgumentIsMissing
	}

	encryptedPassword, err := us.crypto.encryptData(user.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	user.Password = []byte(base64.StdEncoding.EncodeToString(encryptedPassword))

	token, salt, err := us.Client.SignUpUser(ctx, user)
	if err != nil {
		return fmt.Errorf("server failed to sign up user: %w", err)
	}

	if err = us.cnfg.SetSalt([]byte(salt)); err != nil {
		return err
	}

	if err := us.Client.SetJWTToken(token); err != nil {
		return err
	}

	return nil
}

func (us *UserService) SignInUser(ctx context.Context, user *models.User) error {
	if user.Login == "" || user.Password == nil {
		return errs.ErrRequiredArgumentIsMissing
	}

	encryptedPassword, err := us.crypto.encryptData(user.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	user.Password = []byte(base64.StdEncoding.EncodeToString(encryptedPassword))

	token, salt, err := us.Client.SignInUser(ctx, user)
	if err != nil {
		return fmt.Errorf("server failed to sign in user: %w", err)
	}
	err = us.crypto.setSalt(salt)
	if err != nil {
		return err
	}

	if err := us.Client.SetJWTToken(token); err != nil {
		return err
	}

	return nil
}

func (us *UserService) SetMasterKey(masterPassword string) error {
	if err := us.cnfg.SetMasterPassword(masterPassword); err != nil {
		return err
	}
	masterKey, err := us.crypto.generateMasterKey()
	if err != nil {
		return err
	}
	return us.cnfg.SetMasterKey(masterKey)
}

func (us *UserService) Logout() error {
	if err := us.cnfg.SetMasterKey(nil); err != nil {
		return err
	}
	if err := us.cnfg.SetMasterPassword(""); err != nil {
		return err
	}
	if err := us.cnfg.SetSalt(nil); err != nil {
		return err
	}
	return nil
}
