package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type response struct {
		Error       string `json:"error,omitempty"`
		CleanedBody string `json:"cleaned_body,omitempty"`
	}
	w.Header().Set("Content-Type", "application/json")

	// Decode input
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters: %s", err)
		resp500 := response{Error: "Something went wrong"}
		dat500, _ := json.Marshal(resp500)
		w.WriteHeader(500)
		w.Write(dat500)
		return
	}

	// Validate chirp length
	if len(params.Body) > 140 {
		resp := response{Error: "Chirp is too long"}
		dat, _ := json.Marshal(resp)
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
	// Valid chirp
	resp := response{CleanedBody: cleaned}
	dat, err := json.Marshal(resp)
	if err != nil {
		log.Printf("error marshalling response: %s", err)
		resp500 := response{Error: "Something went wrong"}
		dat500, _ := json.Marshal(resp500)
		w.WriteHeader(500)
		w.Write(dat500)
		return
	}

	// Send valid response
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
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
