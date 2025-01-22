package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/michalronin/chirpy/internal/auth"
	"github.com/michalronin/chirpy/internal/database"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("bearer token not found")
		w.WriteHeader(401)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %s", err)
		w.WriteHeader(401)
		return
	}
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		log.Printf("error parsing chirp ID: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if chirp.UserID != userID {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if err := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpUUID,
		UserID: userID,
	}); err != nil {
		log.Printf("chirp not found: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
