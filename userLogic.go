package main

import (
	_ "encoding/json"
	"net/http"
)

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) postUser(w http.ResponseWriter, r *http.Request) {

}

// Copy Write Get Chirps Logic for user, problably add users/chirp in the other just load data and write the database/add to DBStructure before writing
