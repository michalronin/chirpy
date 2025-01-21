package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/michalronin/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("incorrect email or password")
		w.WriteHeader(401)
		return
	}
	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		log.Printf("incorrect email or password")
		w.WriteHeader(401)
		return
	}
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = 3600
	} else if params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = 3600
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		log.Printf("error generating token: %s", err)
		w.WriteHeader(401)
		return
	}
	resp := response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}
	encoded, err := json.Marshal(resp)
	if err != nil {
		log.Printf("incorrect email or password")
		w.WriteHeader(401)
		return
	}
	w.WriteHeader(200)
	w.Write(encoded)
}
