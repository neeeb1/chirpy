package api

import (
	"net/http"
	"time"

	"github.com/neeeb1/chirpy/internal/auth"
	"github.com/neeeb1/chirpy/internal/database"
)

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	dbToken, err := cfg.DbQueries.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	if dbToken.ExpiresAt.Before(time.Now()) || dbToken.RevokedAt.Valid {
		respondWithError(w, 401, "token expired")
		return
	}

	userID, err := cfg.DbQueries.GetUserIDByRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	newJWT, err := auth.MakeJWT(userID.UUID, cfg.Secret, (time.Hour))
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	params := database.CreateRefreshTokenParams{
		Token:     newJWT,
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	_, err = cfg.DbQueries.CreateRefreshToken(r.Context(), params)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: newJWT,
	}

	respondWithJSON(w, 200, resp)

}

func (cfg *ApiConfig) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	_, err = cfg.DbQueries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	respondWithJSON(w, 204, "")

}
