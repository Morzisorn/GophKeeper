package hash

import "golang.org/x/crypto/bcrypt"

func GetHash(body []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(body, bcrypt.DefaultCost)
}

func VerifyHash(password, hash []byte) bool {
    err := bcrypt.CompareHashAndPassword(hash, password)
    return err == nil
}