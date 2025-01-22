package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/michalronin/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	sorting := r.URL.Query().Get("sort")
	chirps := []database.Chirp{}
	if len(s) == 0 {
		allChirps, err := cfg.db.GetAllChirps(r.Context())
		if err != nil {
			fmt.Printf("error retrieving chirps: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			return
		}
		chirps = allChirps
	} else {
		userID, err := uuid.Parse(s)
		if err != nil {
			fmt.Printf("error parsing user ID: %s", err)
			w.WriteHeader(500)
			return
		}
		chirpsPerUser, err := cfg.db.GetAllChirpsForUser(r.Context(), userID)
		if err != nil {
			fmt.Printf("error retrieving chirps: %s", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(500)
			return
		}
		chirps = chirpsPerUser
	}

	resp := []Chirp{}
	if len(sorting) == 0 || sorting == "asc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	} else if sorting == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.After(chirps[j].CreatedAt) })
	}
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
