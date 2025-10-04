package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/neeeb1/chirpy/internal/auth"
)

func (cfg *ApiConfig) HandlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	b := body{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&b)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	keyString, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	if keyString != cfg.PolkaKey {
		respondWithError(w, 401, "")
		return
	}

	if b.Event != "user.upgraded" {
		respondWithJSON(w, 204, "")
		return
	}

	_, err = cfg.DbQueries.UpgradeChirpyRed(r.Context(), uuid.MustParse(b.Data.UserID))
	if err != nil {
		respondWithError(w, 404, err.Error())
		return
	}

	respondWithJSON(w, 204, "")

}
