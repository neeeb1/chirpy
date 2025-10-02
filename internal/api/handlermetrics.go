package api

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, req *http.Request) {

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

func (cfg *ApiConfig) HandlerResetMetrics(w http.ResponseWriter, req *http.Request) {
	cfg.serverHits.Swap(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("Hits reset to %d", cfg.serverHits.Load())

	w.Write([]byte(body))
}

func (cfg *ApiConfig) MiddlewareMetricsIncr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.serverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
