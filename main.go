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
	handler := http.FileServer(http.Dir("./"))
	mux.Handle("/", handler)
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
