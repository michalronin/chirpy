package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("token not found in request header")
	}
	apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
	return apiKey, nil
}
