package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lib/pq"
)

func IsDuplicatedKeys(err error) bool {
	var pqErr *pq.Error
	// duplicated keys error has code 23505
	return err != nil && errors.As(err, &pqErr) && pqErr.Code == "23505"
}

func respondWithError(
	w http.ResponseWriter,
	code int,
	message string,
) {
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

func respondWithJSON(
	w http.ResponseWriter,
	code int,
	payload any,
) {
	w.Header().Add("ContentType", "application/json")
	w.WriteHeader(code)
	data, err := json.Marshal(payload)
	if err != nil {
		// TODO: handle error
		return
	}
	w.Write(data)
}
