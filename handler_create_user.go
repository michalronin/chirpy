package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	w.Header().Set("Content-Type", "application/json")

	// decode
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("error decoding parameters: %s", err)
		return
	}
	if len(params.Email) == 0 {
		w.WriteHeader(400)
		log.Printf("error: email not provided: %s", params)
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("error creating user: %s", err)
		return
	}
	resp := response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	encoded, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("error encoding response: %s", err)
		return
	}
	w.WriteHeader(201)
	w.Write(encoded)
}
