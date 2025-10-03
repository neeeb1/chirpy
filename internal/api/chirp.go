package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/neeeb1/chirpy/internal/database"
)

type chirp struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

func (cfg *ApiConfig) HandlerChirp(w http.ResponseWriter, r *http.Request) {
	var c chirp
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&c)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	validChirp, err := validateChirp(c)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}

	userId, err := uuid.Parse(validChirp.UserID)
	if err != nil {
		respondWithError(w, 400, "failed to parse uuid")
	}
	params := database.CreateChirpParams{
		Body:   validChirp.Body,
		UserID: uuid.NullUUID{userId, true},
	}

	newChirp, err := cfg.DbQuereies.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, 400, "failed to create chirp")
	}

	type response struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	res := response{
		ID:        string(newChirp.ID.String()),
		CreatedAt: newChirp.CreatedAt.String(),
		UpdatedAt: newChirp.UpdatedAt.String(),
		Body:      newChirp.Body,
		UserID:    newChirp.UserID.UUID.String(),
	}

	respondWithJSON(w, 201, res)
}

func validateChirp(c chirp) (chirp, error) {
	validChirp := chirp{}
	if len(c.Body) <= 140 {
		validChirp = profanityFilter(c)
		return validChirp, nil
	} else {
		return validChirp, fmt.Errorf("chirp is too long")
	}
}

func profanityFilter(c chirp) chirp {
	words := strings.Split(c.Body, " ")
	for i, w := range words {
		switch strings.ToLower(w) {
		case "kerfuffle":
			words[i] = "****"
		case "sharbert":
			words[i] = "****"
		case "fornax":
			words[i] = "****"
		default:
			continue
		}
	}

	cleanBody := strings.Join(words, " ")

	cleanChrip := chirp{
		Body:   cleanBody,
		UserID: c.UserID,
	}

	return cleanChrip
}
