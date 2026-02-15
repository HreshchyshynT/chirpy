package main

import (
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()

	var root http.Dir

	serveMux.Handle("/", http.FileServer(root))

	var server http.Server
	server.Addr = ":8080"
	server.Handler = serveMux
	server.ListenAndServe()
}
