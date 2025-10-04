package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/neeeb1/chirpy/internal/auth"
	"github.com/neeeb1/chirpy/internal/database"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *ApiConfig) HandlerNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type body struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	b := body{}

	err := decoder.Decode(&b)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to decode request: %s", err))
		return
	}

	hash, err := auth.HashPassword(b.Password)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to hash password: %s", err))
		return
	}

	params := database.CreateUserParams{
		HashedPassword: hash,
		Email:          b.Email,
	}

	dbUser, err := cfg.DbQueries.CreateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to create user: %s", err))
		return
	}

	newUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, 201, newUser)
}

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type body struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	b := body{}

	err := decoder.Decode(&b)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to decode request: %s", err))
		return
	}

	user, err := cfg.DbQueries.GetUserByEmail(r.Context(), b.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(b.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	if match {
		token, err := auth.MakeJWT(user.ID, cfg.Secret, (time.Hour))
		if err != nil {
			respondWithError(w, 401, "Failed to generate user token")
			return
		}

		refreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, 401, "Failed to generate refresh token")
			return
		}

		refreshTokenParams := database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		}

		_, err = cfg.DbQueries.CreateRefreshToken(r.Context(), refreshTokenParams)
		if err != nil {
			respondWithError(w, 401, "Failed to create refresh token")
			return
		}

		resp := User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refreshToken,
		}
		respondWithJSON(w, 200, resp)
		return
	} else {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

}

func (cfg *ApiConfig) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type body struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	b := body{}

	err := decoder.Decode(&b)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("failed to decode request: %s", err))
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	hash, err := auth.HashPassword(b.Password)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("failed to hash password: %s", err))
		return
	}

	params := database.UpdateUserEmailAndPasswordParams{
		Email:          b.Email,
		HashedPassword: hash,
		ID:             userID,
	}

	updatedUser, err := cfg.DbQueries.UpdateUserEmailAndPassword(r.Context(), params)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	u := User{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}

	respondWithJSON(w, 200, u)

}
