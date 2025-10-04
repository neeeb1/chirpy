package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	str, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return str, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var userID uuid.UUID

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return userID, err
	}

	user, err := token.Claims.GetSubject()
	if err != nil {
		return userID, err
	}

	userID, err = uuid.Parse(user)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

func MakeRefreshToken() (string, error) {
	var token string

	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return token, err
	}

	token = hex.EncodeToString(data)
	return token, nil
}
