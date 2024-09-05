package main

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var ErrNoUser = errors.New("User not found")

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserReturn struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) postUser(w http.ResponseWriter, r *http.Request) {
	validateUser(w, r, db)
	db.id++
}

func validateUser(w http.ResponseWriter, r *http.Request, db *DB) {
	const cost = 15
	decoder := json.NewDecoder(r.Body)
	user := User{
		Id:       db.id,
		Email:    "",
		Password: "",
	}
	err := decoder.Decode(&user)
	if err != nil {
		log.Printf("Error decoding user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding user")
		return
	}

	hashedPW, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Error PW")
		return
	}
	user.Password = string(hashedPW)

	userR, err := db.CreateUser(user)
	if err != nil {
		log.Printf("Error Creating and Saving User: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Database Error")
		return
	}
	respondWithJson(w, http.StatusCreated, userR)
	return
}

func (db *DB) getUserByMail(email string) (User, error) {
	dbs, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	for _, user := range dbs.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrNoUser
}

func (db *DB) validateMail(email string) (bool, error) {
	_, err := db.getUserByMail(email)
	if err == ErrNoUser {
		return true, nil
	}
	return false, err
}

// Copy Write Get Chirps Logic for user, problably add users/chirp in the other just load data and write the database/add to DBStructure before writing
