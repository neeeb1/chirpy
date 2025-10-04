package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/neeeb1/chirpy/internal/auth"
	"github.com/neeeb1/chirpy/internal/database"
)

type chirp struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

func (cfg *ApiConfig) HandlerPostChirp(w http.ResponseWriter, r *http.Request) {
	var c chirp
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&c)
	if err != nil {
		respondWithError(w, 401, err.Error())
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

	c.UserID = userID.String()

	validChirp, err := validateChirp(c)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	userId, err := uuid.Parse(validChirp.UserID)
	if err != nil {
		respondWithError(w, 401, "failed to parse uuid")
		return
	}
	params := database.CreateChirpParams{
		Body:   validChirp.Body,
		UserID: uuid.NullUUID{UUID: userId, Valid: true},
	}

	newChirp, err := cfg.DbQueries.CreateChirp(r.Context(), params)
	if err != nil {
		respondWithError(w, 401, "failed to create chirp")
		return
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

func (cfg *ApiConfig) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("author_id")

	var dbChirps []database.Chirp
	var err error

	if userID != "" {
		dbChirps, err = cfg.DbQueries.GetChirpByUserID(r.Context(), uuid.NullUUID{UUID: uuid.MustParse(userID), Valid: true})
		if err != nil || len(dbChirps) == 0 {
			respondWithError(w, 400, "failed to get chirps")
			return
		}
	} else {
		dbChirps, err = cfg.DbQueries.GetAllChirps(r.Context())
		if err != nil || len(dbChirps) == 0 {
			respondWithError(w, 400, "failed to get chirps")
			return
		}
	}

	type Chirps struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	res := make([]Chirps, len(dbChirps))

	for i, c := range dbChirps {
		res[i] = Chirps{
			ID:        c.ID.String(),
			CreatedAt: c.CreatedAt.String(),
			UpdatedAt: c.UpdatedAt.String(),
			Body:      c.Body,
			UserID:    c.UserID.UUID.String(),
		}
	}

	respondWithJSON(w, 200, res)
}

func (cfg *ApiConfig) HandlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 400, "failed to parse chirp id")
		return
	}

	chirp, err := cfg.DbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "chirp not found")
		return
	}

	type response struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Body      string `json:"body"`
		UserID    string `json:"user_id"`
	}

	res := response{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt.String(),
		UpdatedAt: chirp.UpdatedAt.String(),
		Body:      chirp.Body,
		UserID:    chirp.UserID.UUID.String(),
	}
	respondWithJSON(w, 200, res)
}

func (cfg *ApiConfig) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 401, "failed to parse chirp id")
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

	chirp, err := cfg.DbQueries.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "chirp not found")
		return
	}

	if chirp.UserID.UUID == userID {
		cfg.DbQueries.DeleteChripByID(r.Context(), chirpID)
		respondWithJSON(w, 204, "")
	} else {
		respondWithJSON(w, 403, "Unauthorized")
	}
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
	var cleanBody []string
	for _, w := range words {
		switch strings.ToLower(w) {
		case "kerfuffle":
			cleanBody = append(cleanBody, "****")
		case "sharbert":
			cleanBody = append(cleanBody, "****")
		case "fornax":
			cleanBody = append(cleanBody, "****")
		default:
			cleanBody = append(cleanBody, w)
		}
	}

	cleanChrip := chirp{
		Body:   strings.Join(cleanBody, " "),
		UserID: c.UserID,
	}

	return cleanChrip
}
