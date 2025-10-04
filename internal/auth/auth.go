package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return tokenString, fmt.Errorf("failed to get authorization header")
	}

	split := strings.Split(tokenString, " ")

	return split[1], nil
}

func GetAPIKey(headers http.Header) (string, error) {
	keyString := headers.Get("Authorization")
	if keyString == "" {
		return keyString, fmt.Errorf("failed to get authorization header")
	}

	split := strings.Split(keyString, " ")

	return split[1], nil
}
