package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/michalronin/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error reading bearer token: %s", err)
		w.WriteHeader(401)
		return
	}
	dbToken, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("bearer token not found in database: %s", err)
		w.WriteHeader(401)
		return
	}
	newToken, err := auth.MakeJWT(dbToken.UserID, cfg.secret)
	if err != nil {
		log.Printf("access token generation failed: %s", err)
		w.WriteHeader(401)
		return
	}
	resp := response{
		Token: newToken,
	}
	encoded, err := json.Marshal(resp)
	if err != nil {
		log.Printf("token encoding failed: %s", err)
		w.WriteHeader(401)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(encoded)
}
