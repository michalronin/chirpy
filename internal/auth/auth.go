package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("token not found in request header")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", errors.New("failed to generate random bytes")
	}
	encoded := hex.EncodeToString(randomBytes)
	return encoded, nil
}
