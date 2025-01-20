package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		fmt.Printf("error parsing chirp ID: %s", err)
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		fmt.Printf("error: chirp with given id not found; %s", err)
		w.WriteHeader(404)
		return
	}
	resp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	encoded, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error encoding json: %s", err)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
	w.Write(encoded)
}
