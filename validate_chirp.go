package main

import (
	"encoding/json"
	"net/http"
)

type chirp struct {
	Body string `json:"body"`
}

type error struct {
	Error string `json:"error"`
}

type valid struct {
	Valid bool `json:"valid"`
}

func (cfg *apiConfig) validateChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	chirp := chirp{}

	err := decoder.Decode(&chirp)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		parseError := error{Error: "Something went wrong"}
		data, _ := json.Marshal(parseError)
		w.Write(data)
		return
	}

	if len(chirp.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		parseError := error{Error: "Chirp is too long"}
		data, _ := json.Marshal(parseError)
		w.Write(data)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	valid := valid{Valid: true}
	data, _ := json.Marshal(valid)
	w.Write(data)

}
