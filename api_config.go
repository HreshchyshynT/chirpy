package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/auth"
	"github.com/hreshchyshynt/chirpy/internal/database"
)

type dbHandler func(
	w http.ResponseWriter,
	r *http.Request,
	queries *database.Queries,
)

type apiHandler func(
	w http.ResponseWriter,
	r *http.Request,
	apiConfig *ApiConfig,
)

type Platform string

const (
	dev Platform = "dev"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Queries        *database.Queries
	Platform       Platform
	Secret         string
}

func (ac *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ac.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (ac *ApiConfig) middlewareDbAccess(next dbHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r, ac.Queries)
	})
}

func (ac *ApiConfig) middlewareWithConfig(next apiHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r, ac)
	})
}

func (ac *ApiConfig) middlewareRequireJWT(
	next func(w http.ResponseWriter, r *http.Request, userId uuid.UUID),
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "No token provided", err)
			return
		}

		id, err := auth.ValidateJWT(token, ac.Secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
			return
		}

		next(w, r, id)
	})
}

func (ac *ApiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, template, ac.FileserverHits.Load())
}

func (ac *ApiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if ac.Platform != dev {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	ac.FileserverHits.Store(0)
	ac.Queries.ClearUsers(r.Context())
	w.WriteHeader(http.StatusOK)
}
