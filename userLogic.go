package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) postUser(w http.ResponseWriter, r *http.Request) {
	validateUser(w, r, db)
	db.id++
}

func validateUser(w http.ResponseWriter, r *http.Request, db *DB) {
	decoder := json.NewDecoder(r.Body)
	user := User{
		Id:    db.id,
		Email: "",
	}
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Error decoding user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding user")
		return
	}
	_, err = db.CreateUser(user)
	if err != nil {
		log.Printf("Error Creating and Saving User: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Database Error")
		return
	}
	respondWithJson(w, http.StatusCreated, user)
	return
}

// Copy Write Get Chirps Logic for user, problably add users/chirp in the other just load data and write the database/add to DBStructure before writing
