package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
	id   int
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]Chirp `json:"users"`
}

func NewDB(path string) (*DB, error) {
	return &DB{
		path: path,
		mux:  &sync.RWMutex{},
		id:   1,
	}, nil
}

func (db *DB) CreateChirp(chirp Chirp) (Chirp, error) {
	err := db.ensureDB()
	chirpSet, err := db.GetChirps()
	if err != nil {
		return Chirp{}, err
	}

	chirpSet = append(chirpSet, chirp)

	err = db.writeChirps(chirpSet)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	chirpMap, err := db.loadDB()
	if err != nil {
		log.Printf("loadDB err")
		return []Chirp{}, err
	}
	chirpSet := []Chirp{}
	var keys = []int{}
	for key := range chirpMap.Chirps {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, k := range keys {
		chirpSet = append(chirpSet, chirpMap.Chirps[k])
	}

	return chirpSet, nil
}

func (db *DB) writeChirps(chirpSet []Chirp) error {
	chirpMap := make(map[int]Chirp)
	for _, chirp := range chirpSet {
		chirpMap[chirp.Id] = chirp
	}
	err := db.writeDB(DBStructure{Chirps: chirpMap})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	dbs := DBStructure{}
	chirpJson, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	err = json.Unmarshal(chirpJson, &dbs)
	if err != nil {
		return DBStructure{}, err
	}
	return dbs, nil
}

func (db *DB) writeDB(dbs DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	valid, err := json.Marshal(dbs)
	if err != nil {
		return err
	}
	err = os.WriteFile(db.path, valid, 0666)
	if err != nil {
		return err
	}
	return nil

}

func (db *DB) ensureDB() error {
	const chirpMap = `
	[
		{
			"chirps": {}
		},
		{
			"user": {}
		}
	]`
	_, err := os.ReadFile(db.path)
	if err != nil {
		err = os.WriteFile(db.path, []byte(chirpMap), 0666)
	}
	if err != nil {
		return err
	}
	return nil
}
