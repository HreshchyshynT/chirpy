package main

import (
	"net/http"
	"strings"
)

func main() {
	serveMux := http.NewServeMux()

	var root http.Dir

	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(root)))
	serveMux.HandleFunc("/healthz", checkHealth)

	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}
	server.ListenAndServe()
}

func checkHealth(rw http.ResponseWriter, r *http.Request) {
	var builder strings.Builder
	builder.WriteString("ContentType: text/plain; charset=utf-8")
	rw.Header().Write(&builder)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}
