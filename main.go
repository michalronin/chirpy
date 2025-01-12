package main

import (
	"log"
	"net/http"
)

func main() {
	filePathRoot := "./"
	port := "8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	mux.Handle("/", http.FileServer(http.Dir("./")))
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
