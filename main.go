package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", HandlerHealth)
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

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
