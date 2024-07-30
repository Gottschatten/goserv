package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	respondWithJSON(w, code, msg)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func (db *dbCounter) validateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	chirp := chirp{}
	const maxBodyLength = 140

	err := decoder.Decode(&chirp)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid request payload")
		return
	}

	if len(chirp.Body) > maxBodyLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	db.ChirpCount++
	chirp.Body = checkProfanity(chirp.Body)
	chirp.Id = db.ChirpCount
	saveChirp(db.DbAdress, &chirp)

	log.Printf("Chirp is valid: %s", chirp.Body)

	respondWithJSON(w, http.StatusOK, chirp)

}

//func postChirp(w http.ResponseWriter, r *http.Request) {
//	db.validateChirp(w, r)
//}

var profanities = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

func checkProfanity(body string) string {
	bodySplit := strings.Split(body, " ")
	bodySplitLower := strings.Split(strings.ToLower(body), " ")
	for i, word := range bodySplitLower {
		if profanities[word] {
			bodySplit[i] = "****"
		}
	}
	body = strings.Join(bodySplit, " ")
	return body
}
