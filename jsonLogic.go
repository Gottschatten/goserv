package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Chirp struct {
	Body string `json:"body"`
}

type returnChirp struct {
	CleanedBody string `json:"cleaned_body"`
}

type Invalid struct {
	Error string `json:"error"`
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	//
	const chirplen = 140
	var badWords = map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	//
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding chirp")
		return
	}
	if len(chirp.Body) > chirplen {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	cleaned := cleanChirp(chirp.Body, badWords)
	respondWithJson(w, http.StatusOK, returnChirp{CleanedBody: cleaned})
	return
}

func cleanChirp(body string, badWords map[string]bool) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if badWords[lowerWord] {
			words[i] = "****"
		}
	}
	joined := strings.Join(words, " ")
	return joined

}

func respondWithError(w http.ResponseWriter, statuscode int, msg string) {
	if statuscode > 499 {
		log.Printf("Serverside error: %s", msg)
	}
	respondWithJson(w, statuscode, Invalid{Error: msg})
}

func respondWithJson(w http.ResponseWriter, statuscode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	valid, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to Marshal json: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statuscode)
	w.Write(valid)

}
