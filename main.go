package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries        *database.Queries
}

func (ac *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ac.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can not open database connection")
	}

	serveMux := http.NewServeMux()

	var root http.Dir
	config := apiConfig{
		queries: database.New(db),
	}

	serveMux.Handle(
		"/app/",
		http.StripPrefix("/app", config.middlewareMetricsInc(http.FileServer(root))),
	)
	serveMux.HandleFunc("GET /api/healthz", checkHealth)
	serveMux.HandleFunc("POST /api/validate_chirp", validateChirp)
	serveMux.HandleFunc("GET /admin/metrics", config.handleMetrics)
	serveMux.HandleFunc("POST /admin/reset", config.handleReset)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	var builder strings.Builder
	builder.WriteString("ContentType: text/plain; charset=utf-8")
	w.Header().Write(&builder)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (ac *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	var builder strings.Builder
	template := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	builder.WriteString("ContentType: text/html")
	w.Header().Write(&builder)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, template, ac.fileserverHits.Load())
}

func (ac *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	ac.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
