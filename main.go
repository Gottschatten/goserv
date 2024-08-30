package main

import (
	"log"
	"net/http"
)

func main() {
	const dir = http.Dir(".")
	const port = "8080"
	const path = "./database.json"

	//
	cfg := apiConfig{
		fileserverHits: 0,
	}

	db, err := NewDB(path)
	if err != nil {
		log.Printf("Error connecting DB: %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		cfg.mwMetricsInc(http.StripPrefix("/app", http.FileServer(dir))),
	)

	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /admin/metrics", cfg.adminMetric)
	mux.HandleFunc("/api/reset", cfg.resetMetrics)
	mux.HandleFunc("POST /api/chirps", db.postChirp)
	mux.HandleFunc("GET /api/chirps", db.getChirp)
	mux.HandleFunc("GET /api/chirps/{chirpId}", db.getOneChirp)

	// Server Config
	app := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Listening on: localhost%s", app.Addr)
	log.Fatal(app.ListenAndServe())
}
