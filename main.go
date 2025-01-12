package main

import (
	"log"
	"net/http"
)

func main() {
	filePathRoot := "./app"
	port := "8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	handler := http.FileServer(http.Dir("./"))
	mux.Handle("/app/", http.StripPrefix("/app/", handler))
	mux.HandleFunc("/healthz", readinessHandler)
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
