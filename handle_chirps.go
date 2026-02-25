package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func toDomainChirp(c database.Chirp) Chirp {
	return Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	}
}

func handleCreateChirp(
	w http.ResponseWriter,
	r *http.Request,
	db *database.Queries,
) {
	type requestBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	defer r.Body.Close()

	var request requestBody

	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, messageInvalidRequestBody, err)
		return
	}

	cleanedText, err := validateChirp(request.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedText,
		UserID: request.UserID,
	})
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Sprintf("Can't save chirp: %v\n", err),
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusCreated, toDomainChirp(chirp))

}
