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
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *ApiConfig) HandlerNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type user struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	u := user{}

	err := decoder.Decode(&u)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to decode request: %s", err))
		return
	}

	hash, err := auth.HashPassword(u.Password)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to hash password: %s", err))
		return
	}

	params := database.CreateUserParams{
		HashedPassword: hash,
		Email:          u.Email,
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

	type user struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	u := user{}

	err := decoder.Decode(&u)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to decode request: %s", err))
		return
	}

	result, err := cfg.DbQueries.GetUserByEmail(r.Context(), u.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(u.Password, result.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	if match {
		resp := User{
			ID:        result.ID,
			CreatedAt: result.CreatedAt,
			UpdatedAt: result.UpdatedAt,
			Email:     result.Email,
		}
		respondWithJSON(w, 200, resp)
		return
	} else {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

}
