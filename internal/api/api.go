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
	DbQueries  *database.Queries
	Platform   string
	Secret     string
}

func RegisterEndpoints(mux *http.ServeMux, apiCfg *ApiConfig) {
	mux.HandleFunc("GET /api/healthz", HandlerHealth)

	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)

	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerPostChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.HandlerGetChirpByID)

	mux.HandleFunc("POST /api/users", apiCfg.HandlerNewUser)
	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.Handle("/app/", apiCfg.MiddlewareMetricsIncr(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
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
