package user_service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"gophkeeper/config"
	"gophkeeper/internal/errs"
	"gophkeeper/internal/hash"
	"gophkeeper/internal/server/repositories"
	"gophkeeper/models"

	"github.com/jackc/pgx/v5"
)

type UserService struct {
	repo repositories.Storage
}

func NewUserService(repo repositories.Storage) *UserService {
	return &UserService{repo: repo}
}

func (us *UserService) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	var err error
	user, err = us.repo.GetUser(ctx, user.Login)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil, fmt.Errorf("get user error: %w", errs.ErrUserNotFound)
	case err == nil:
		return user, nil
	}
	return nil, fmt.Errorf("get user error: %w", err)
}

func (us *UserService) SignUpUser(ctx context.Context, login, encryptedPassword string) (token string, err error) {
	_, err = us.GetUser(ctx, &models.User{Login: login})
	switch {
	case err == nil:
		return "", fmt.Errorf("sign up user error: %w", errs.ErrUserAlreadyRegistered)
	case !errors.Is(err, errs.ErrUserNotFound):
		return "", fmt.Errorf("sign up user error: %w", err)
	}

	decryptedPassword, err := decryptPassword(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt password: %w", err)
	}

	hash, err := hash.GetHash(decryptedPassword)
	if err != nil {
		return "", fmt.Errorf("generate password hash error: %w", err)
	}
	err = us.repo.SignUpUser(ctx, &models.User{
		Login:    login,
		Password: hash,
	})
	if err != nil {
		return "", fmt.Errorf("sign up user in db error: %w", err)
	}

	token, err = generateToken(login)
	if err != nil {
		return "", fmt.Errorf("generate token after sign up user error: %w", err)
	}

	return token, nil
}

func (us *UserService) SignInUser(ctx context.Context, login, encryptedPassword string) (token string, err error) {
	user, err := us.GetUser(ctx, &models.User{Login: login})
	switch {
	case errors.Is(err, errs.ErrUserNotFound):
		return "", err
	case err != nil && !errors.Is(err, errs.ErrUserNotFound):
		return "", fmt.Errorf("sign in user error: %w", err)
	}

	decryptedPassword, err := decryptPassword(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt password: %w", err)
	}

	if !hash.VerifyHash(decryptedPassword, user.Password) {
		return "", errs.ErrIncorrectCredentials
	}

	token, err = generateToken(login)
	if err != nil {
		return "", fmt.Errorf("generate token after sign up user error: %w", err)
	}

	return token, nil
}

func decryptPassword(encryptedPassword string) ([]byte, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return nil, fmt.Errorf("decode password from base64 error: %v", err)
	}

	decryptedPassword, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		config.GetServerConfig().PrivateKey,
		encryptedBytes,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	return decryptedPassword, nil
}
