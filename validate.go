package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var bodyDecoded requestBody
	err := decoder.Decode(&bodyDecoded)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding request body")
		return
	}
	if len(bodyDecoded.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := cleanText(bodyDecoded.Body)

	body := responseBody{
		CleanedBody: cleanedBody,
	}
	respondWithJSON(w, http.StatusOK, body)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	type errorBody struct {
		Message string `json:"error"`
	}

	w.Header().Add("ContentType", "application/json")
	w.WriteHeader(code)

	body := errorBody{
		Message: message,
	}

	data, err := json.Marshal(body)

	if err != nil {
		// TODO: handle error
		return
	}
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("ContentType", "application/json")
	w.WriteHeader(code)
	data, err := json.Marshal(payload)
	if err != nil {
		// TODO: handle error
		return
	}
	w.Write(data)
}

func cleanText(input string) string {
	var builder strings.Builder

	profanes := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	splitted := strings.Split(input, " ")

	for i, s := range splitted {

		if _, ok := profanes[strings.ToLower(s)]; ok {
			builder.WriteString("****")
		} else {
			builder.WriteString(s)
		}

		if i < len(splitted)-1 {
			builder.WriteRune(' ')
		}
	}

	return builder.String()
}
