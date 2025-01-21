package main

import (
	"log"
	"net/http"

	"github.com/michalronin/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("bearer token not found: %s", err)
		w.WriteHeader(401)
		return
	}
	if err := cfg.db.RevokeRefreshToken(r.Context(), token); err != nil {
		log.Printf("error revoking token: %s", err)
		w.WriteHeader(401)
		return
	}
	w.WriteHeader(204)
}
