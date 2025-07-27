package user_service

import (
	"fmt"
	"time"

	"gophkeeper/config"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken(login string) (string, error) {
	claims := jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(7 * time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	cnfg, err := config.GetServerConfig()
	if err != nil {
		return "", fmt.Errorf("get server config error: %w", err)
	}
	signedToken, err := token.SignedString([]byte(cnfg.SecretKey))
	if err != nil {
		return "", fmt.Errorf("generate token error: %w", err)
	}

	return signedToken, nil
}
