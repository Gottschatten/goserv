package main

import (
	"log"
	"net/http"
)

func main() {
	const dir = http.Dir(".")
	cfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		cfg.mwMetricsInc(http.StripPrefix("/app", http.FileServer(dir))),
	)

	mux.HandleFunc("GET /api/healthz", healthz)
	mux.HandleFunc("GET /api/metrics", cfg.fileserverMetrics)
	mux.HandleFunc("/api/reset", cfg.resetMetrics)

	// Server Config
	app := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(app.ListenAndServe())
}
