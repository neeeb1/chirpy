package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

func (cfg *ApiConfig) HandlerValidater(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	c := chirp{}

	err := decoder.Decode(&c)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error decoding chirp: %s", err))
	}

	if len(c.Body) <= 140 {
		type cleanedChirp struct {
			CleanedBody string `json:"cleaned_body"`
		}

		respBody := cleanedChirp{
			CleanedBody: profanityFilter(c),
		}

		respondWithJSON(w, 200, respBody)
	} else {
		respondWithError(w, 400, "Chirp is too long")
	}

}

func profanityFilter(c chirp) string {

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
	cleaned := strings.Join(words, " ")
	return cleaned
}
