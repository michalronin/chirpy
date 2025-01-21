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

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
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
	token, err := auth.MakeJWT(user.ID, cfg.secret)
	if err != nil {
		log.Printf("error generating token: %s", err)
		w.WriteHeader(401)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()

	if err != nil {
		log.Printf("error generating refresh token: %s", err)
		w.WriteHeader(401)
		return
	}

	if err := cfg.db.SaveRefreshToken(r.Context(), database.SaveRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60), // 60 days
	}); err != nil {
		log.Printf("error storing refresh token: %s", err)
		w.WriteHeader(401)
		return
	}
	resp := response{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
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
