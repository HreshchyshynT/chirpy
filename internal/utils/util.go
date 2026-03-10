package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	MessageInvalidRequestBody  = "Invalid request body"
	MessageInternalServerError = "Internal server error"
)

func RespondWithError(
	w http.ResponseWriter,
	code int,
	message string,
	err error,
) {
	if err != nil {
		log.Println(err)
	}
	type errorBody struct {
		Message string `json:"error"`
	}

	w.Header().Add("ContentType", "application/json")
	w.WriteHeader(code)

	body := errorBody{
		Message: message,
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(&body)

	if err != nil {
		log.Println(err)
	}
}

func RespondWithJSON(
	w http.ResponseWriter,
	code int,
	payload any,
) {
	w.Header().Add("ContentType", "application/json")
	w.WriteHeader(code)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(payload)
}
