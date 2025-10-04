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
