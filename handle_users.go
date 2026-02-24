package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func toDomain(u database.User) User {
	return User{
		ID:        u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email:     u.Email,
	}
}

func handleCreateUser(
	w http.ResponseWriter,
	r *http.Request,
	db *database.Queries,
) {
	type requestBody struct {
		Email string `json:"email"`
	}

	var body requestBody

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := db.CreateUser(r.Context(), body.Email)

	if IsDuplicatedKeys(err) {
		respondWithError(w, http.StatusBadRequest, "Email is already used")
		return
	}

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Can't create user. Try again later",
		)
		return
	}

	respondWithJSON(w, http.StatusCreated, toDomain(user))
}
