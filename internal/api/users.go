package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *ApiConfig) HandlerNewUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type email struct {
		Email string `json:"email"`
	}

	e := email{}

	err := decoder.Decode(&e)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to decode request: %s", err))
	}

	dbUser, err := cfg.DbQuereies.CreateUser(r.Context(), e.Email)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("failed to create user: %s", err))
	}

	newUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, 201, newUser)
}
