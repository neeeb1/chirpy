package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	serverHits atomic.Int32
}

func main() {
	var apiCfg apiConfig
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", HandlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerResetMetrics)
	mux.Handle("/app/", apiCfg.middlewareMetricsIncr(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	server := http.Server{}

	server.Handler = mux
	server.Addr = ":8080"

	server.ListenAndServe()

}

func HandlerHealth(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("status: OK!"))
}

func (cfg *apiConfig) middlewareMetricsIncr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.serverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) HandlerMetrics(w http.ResponseWriter, req *http.Request) {

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.serverHits.Load())

	w.Write([]byte(body))
}

func (cfg *apiConfig) HandlerResetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.serverHits.Swap(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("Hits reset to %d", cfg.serverHits.Load())

	w.Write([]byte(body))
}
