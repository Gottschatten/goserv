package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	cfg := apiConfig{
		fileserverHits: 0,
	}

	const dir = "."
	const port = "8080"
	mux := http.NewServeMux()

	// Serve Webpages, files at /app
	mux.Handle(
		"/app/*",
		cfg.middlewareMetric(http.StripPrefix("/app", http.FileServer(http.Dir(dir)))),
	)

	// admin Endpoints /admin
	mux.HandleFunc("GET /admin/metrics", cfg.sendHits)

	// API Endpoints /api
	mux.HandleFunc("GET /api/healthz", healthzHandle)
	mux.HandleFunc("/api/reset", cfg.resetHits)

	// Define Server adress and Handler
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving ")
	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetric(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset!"))
	log.Printf("Resetting hits")
	cfg.fileserverHits = 0
}

func (cfg *apiConfig) sendHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>
		
</html>`,
		cfg.fileserverHits)))
}

func healthzHandle(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

