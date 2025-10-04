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

	apiCfg.Secret = os.Getenv("SECRET")

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		return
	}
	apiCfg.DbQueries = database.New(db)

	mux := http.NewServeMux()

	api.RegisterEndpoints(mux, &apiCfg)

	server := http.Server{}

	server.Handler = mux
	server.Addr = ":8080"

	server.ListenAndServe()

}
