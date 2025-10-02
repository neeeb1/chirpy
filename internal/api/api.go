package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/neeeb1/chirpy/internal/database"
)

type ApiConfig struct {
	serverHits atomic.Int32
	DbQuereies *database.Queries
	Platform   string
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	respBody := payload

	data, err := json.Marshal(respBody)
	if err != nil {
		fmt.Printf("failed to encode response: %v\n", err)
	}

	w.WriteHeader(code)
	w.Write(data)
	w.Header().Set("Content-Type", "application/json")

}
func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string
	}

	respBody := errorResponse{
		Error: msg,
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		fmt.Printf("failed to encode error response: %v\n", err)
	}

	w.WriteHeader(code)
	w.Write(data)
	w.Header().Set("Content-Type", "application/json")

}
