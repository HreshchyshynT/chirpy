package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/chirps"
	"github.com/hreshchyshynt/chirpy/internal/config"
	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/hreshchyshynt/chirpy/internal/user"
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
	config := config.ApiConfig{
		Queries:  database.New(db),
		Platform: config.Platform(os.Getenv("PLATFORM")),
		Secret:   os.Getenv("SECRET"),
	}

	serveMux.Handle(
		"/app/",
		http.StripPrefix("/app", config.MiddlewareMetricsInc(http.FileServer(root))),
	)
	serveMux.HandleFunc("GET /api/healthz", checkHealth)

	serveMux.Handle(
		"POST /api/chirps",
		config.MiddlewareRequireJWT(
			func(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
				chirps.HandleCreateChirp(w, r, userId, config.Queries)
			},
		),
	)
	serveMux.Handle(
		"DELETE /api/chirps/{chirpID}",
		config.MiddlewareRequireJWT(
			func(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
				chirps.HandleDeleteChirp(w, r, userId, &config)
			},
		),
	)
	serveMux.Handle(
		"GET /api/chirps",
		config.MiddlewareDbAccess(chirps.HandleGetAllChirps),
	)
	serveMux.Handle("GET /api/chirps/{chirpID}",
		config.MiddlewareDbAccess(chirps.HandleGetChirp),
	)

	serveMux.Handle("POST /api/users", config.MiddlewareDbAccess(user.HandleCreateUser))
	serveMux.Handle("POST /api/login", config.MiddlewareWithConfig(user.HandleLogin))
	serveMux.Handle(
		"PUT /api/users",
		config.MiddlewareRequireJWT(
			func(w http.ResponseWriter, r *http.Request, userId uuid.UUID) {
				user.HandleUpdateUser(w, r, userId, &config)
			},
		),
	)

	serveMux.Handle(
		"POST /api/refresh",
		config.MiddlewareWithConfig(user.HandleRefreshToken),
	)
	serveMux.Handle(
		"POST /api/revoke",
		config.MiddlewareWithConfig(user.HandleRevoke),
	)

	serveMux.HandleFunc("GET /admin/metrics", config.HandleMetrics)
	serveMux.HandleFunc("POST /admin/reset", config.HandleReset)

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
