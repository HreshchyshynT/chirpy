package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/hreshchyshynt/chirpy/internal/database"
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
