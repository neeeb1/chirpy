package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neeeb1/chirpy/internal/api"
	"github.com/neeeb1/chirpy/internal/database"
)

func main() {
	godotenv.Load()

	var apiCfg api.ApiConfig

	apiCfg.Platform = os.Getenv("PLATFORM")

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		return
	}
	apiCfg.DbQuereies = database.New(db)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", api.HandlerHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerPostChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerGetChirps)
	mux.HandleFunc("POST /api/users", apiCfg.HandlerNewUser)
	mux.Handle("/app/", apiCfg.MiddlewareMetricsIncr(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	server := http.Server{}

	server.Handler = mux
	server.Addr = ":8080"

	server.ListenAndServe()

}
