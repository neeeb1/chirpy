package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (cfg *ApiConfig) HandlerValidater(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	c := chirp{}
	err := decoder.Decode(&c)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(c.Body) <= 140 {
		type validResponse struct {
			Valid bool `json:"valid"`
		}

		respBody := validResponse{
			Valid: true,
		}

		data, err := json.Marshal(respBody)
		if err != nil {
			fmt.Printf("failed to validate: %v\n", err)
		}

		w.WriteHeader(200)
		w.Write(data)
	} else {
		type errorResponse struct {
			Error string `json:"error"`
		}

		respBody := errorResponse{
			Error: "Chirp is too long",
		}

		data, err := json.Marshal(respBody)
		if err != nil {
			fmt.Printf("failed to validate: %v\n", err)
		}

		w.WriteHeader(400)
		w.Write(data)
	}

	w.Header().Set("Content-Type", "application/json")

}
