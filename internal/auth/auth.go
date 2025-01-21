package auth

import (
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
