package user_service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (us *UserService) generateToken(login string) (string, error) {
	claims := jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(7 * time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(us.cnfg.GetSecretKey()))
	if err != nil {
		return "", fmt.Errorf("generate token error: %w", err)
	}

	return signedToken, nil
}
