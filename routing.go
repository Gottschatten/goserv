package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) middlewareMetric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "test/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset!"))
	log.Printf("Resetting serving hits.")
	cfg.fileserverHits = 0
}
