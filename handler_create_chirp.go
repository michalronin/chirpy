package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/michalronin/chirpy/internal/auth"
	"github.com/michalronin/chirpy/internal/database"
)

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	// decode
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		fmt.Printf("error decoding request body: %s", err)
		w.WriteHeader(400)
		return
	}
	// validate token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error validating token: %s", err)
		w.WriteHeader(401)
		return
	}
	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %s", err)
		w.WriteHeader(401)
		return
	}
	// Validate chirp length
	if len(params.Body) > 140 {
		resp := errorResponse{Error: "Chirp is too long"}
		dat, _ := json.Marshal(resp)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: id,
	})
	if err != nil {
		fmt.Printf("error writing chirp to database: %s", err)
		w.WriteHeader(500)
		return
	}
	resp := response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	encoded, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("error encoding response: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(encoded)
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
