package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	DbAdress string
	mux      sync.Mutex
}

type DBStruct struct {
	Chirps map[int]chirp `json:"chirps"`
}

func deleteDB(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Printf("Error deleting database: %s", err)
	}
}

func saveChirp(path string, chirp *chirp) {
	jChirp, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling chirp: %d", chirp.Id)
		deleteDB(path)
	}
	os.WriteFile(path, jChirp, 0644)
}

func NewDB(path string) (*DB, error) {
	db := &DB{DbAdress: path,
		mux: sync.Mutex{}}

	_, err := os.Create(path)
	if err != nil {
		return &DB{}, err
	}
	return db, nil
}

func (db *DB) loadDB() (DBStruct, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	loadedDB, err := os.ReadFile(db.DbAdress)
	if err != nil {
		return DBStruct{}, err
	}
	var loadedChirps DBStruct

	unmarshellErr := json.Unmarshal(loadedDB, &loadedChirps.Chirps)
	if unmarshellErr != nil {
		return DBStruct{}, err
	}
	return loadedChirps, nil
}

func (db *DB) writeDB(dbs DBStruct) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	f, err := os.OpenFile(db.DbAdress, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	jChirps, err := json.Marshal(dbs)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(jChirps)); err != nil {
		return err
	}
	return nil
}
