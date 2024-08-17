package main

import (
	"log"
	"net/http"
)

func main() {
	const dir = http.Dir(".")
	const port = "8080"

	//
	cfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		cfg.mwMetricsInc(http.StripPrefix("/app", http.FileServer(dir))),
	)

	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /admin/metrics", cfg.adminMetric)
	mux.HandleFunc("/api/reset", cfg.resetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	// Server Config
	app := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Fatal(app.ListenAndServe())
}
