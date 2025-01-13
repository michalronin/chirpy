package main

import (
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

func main() {
	filePathRoot := "./app"
	port := "8080"
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	handler := http.FileServer(http.Dir("./"))
	mux.Handle("/app/", http.StripPrefix("/app/", apiCfg.middlewareMetricsIncrement(handler)))
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /api/reset", apiCfg.resetHandler)
	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	hits := int(cfg.fileserverHits.Load())
	w.Write([]byte("Hits: " + strconv.Itoa(hits)))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = atomic.Int32{}
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsIncrement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
