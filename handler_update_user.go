package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/michalronin/chirpy/internal/auth"
	"github.com/michalronin/chirpy/internal/database"
)

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("bearer token missing: %s", err)
		w.WriteHeader(401)
		return
	}
	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %s", err)
		w.WriteHeader(401)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	decodeErr := decoder.Decode(&params)
	if decodeErr != nil {
		w.WriteHeader(500)
		log.Printf("error decoding parameters: %s", err)
		return
	}
	if len(params.Email) == 0 {
		w.WriteHeader(400)
		log.Printf("error: email not provided: %s", params)
		return
	}
	if len(params.Password) == 0 {
		w.WriteHeader(400)
		log.Printf("error: password is empty: %s", params)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("error hashing password: %s", err)
		return
	}
	updatedUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             id,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		w.WriteHeader(500)
		log.Printf("error updating user: %s", err)
		return
	}
	resp := response{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}
	encoded, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("error encoding response: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(encoded)
}
