package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	passwordBytes := []byte(password)
	hashBytes := []byte(hash)
	if err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes); err != nil {
		return errors.New("checking password hash failed")
	}
	return nil
}
