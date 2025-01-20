package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		fmt.Printf("error retrieving chirps: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		return
	}
	resp := []Chirp{}
	for _, dbChirp := range chirps {
		chirp := Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
		resp = append(resp, chirp)
	}
	encoded, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error encoding chirps: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(encoded)
}
