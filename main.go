package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
		queries:  database.New(db),
		Platform: Platform(os.Getenv("PLATFORM")),
	}

	serveMux.Handle(
		"/app/",
		http.StripPrefix("/app", config.middlewareMetricsInc(http.FileServer(root))),
	)
	serveMux.HandleFunc("GET /api/healthz", checkHealth)

	serveMux.Handle(
		"POST /api/chirps",
		config.middlewareDbAccess(handleCreateChirp),
	)
	serveMux.Handle(
		"GET /api/chirps",
		config.middlewareDbAccess(handleGetAllChirps),
	)

	serveMux.Handle("POST /api/users", config.middlewareDbAccess(handleCreateUser))

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
