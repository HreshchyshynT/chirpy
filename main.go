package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (ac *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ac.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	serveMux := http.NewServeMux()

	var root http.Dir
	var config apiConfig

	serveMux.Handle(
		"/app/",
		http.StripPrefix("/app", config.middlewareMetricsInc(http.FileServer(root))),
	)
	serveMux.HandleFunc("/healthz", checkHealth)
	serveMux.HandleFunc("/metrics", config.handleMetrics)
	serveMux.HandleFunc("/reset", config.handleReset)

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
	builder.WriteString("ContentType: text/plain; charset=utf-8")
	w.Header().Write(&builder)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %v", ac.fileserverHits.Load())
}

func (ac *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	ac.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
