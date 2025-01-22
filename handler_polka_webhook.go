package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/michalronin/chirpy/internal/auth"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	user, err := cfg.db.GetUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("GetUserByID error")
		return
	}
	if err := cfg.db.UpgradeUserToChirpyRed(r.Context(), user.ID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("UpgradeUserToChirpyRed error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
