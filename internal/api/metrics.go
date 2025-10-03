package api

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.serverHits.Load())

	w.Write([]byte(body))
}

func (cfg *ApiConfig) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "forbidden")
		return
	}

	cfg.serverHits.Swap(0)
	cfg.DbQueries.DeleteAllUsers(r.Context())

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf("Hits reset to %d\nRemoved all users from database", cfg.serverHits.Load())

	w.Write([]byte(body))
}

func (cfg *ApiConfig) MiddlewareMetricsIncr(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.serverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
