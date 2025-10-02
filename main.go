package main

import (
	"net/http"

	"github.com/neeeb1/chirpy/internal/api"
)

func main() {
	var apiCfg api.ApiConfig
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", api.HandlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.HandlerValidater)
	mux.Handle("/app/", apiCfg.MiddlewareMetricsIncr(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	server := http.Server{}

	server.Handler = mux
	server.Addr = ":8080"

	server.ListenAndServe()

}
