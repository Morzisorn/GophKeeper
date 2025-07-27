package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateSalt() (string, error) {
	salt := make([]byte, 32)

	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	return base64.StdEncoding.EncodeToString(salt), nil
}
