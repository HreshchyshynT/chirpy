package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/auth"
	"github.com/hreshchyshynt/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func toDomainUser(u database.User) User {
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body requestBody

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, messageInvalidRequestBody, err)
		return
	}

	pass, err := auth.HashPassword(body.Password)
	if err != nil {
		respondWithError(w,
			http.StatusInternalServerError,
			messageInternalServerError,
			err,
		)
		return
	}

	user, err := db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          body.Email,
		HashedPassword: pass,
	})

	if IsDuplicatedKeys(err) {
		respondWithError(w, http.StatusBadRequest, "Email is already used", err)
		return
	}

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Can't create user. Try again later",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusCreated, toDomainUser(user))
}

func handleLogin(
	w http.ResponseWriter,
	r *http.Request,
	db *database.Queries,
) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body requestBody

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	defer r.Body.Close()

	err := decoder.Decode(&body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, messageInvalidRequestBody, err)
		return
	}

	user, err := db.FindUser(r.Context(), body.Email)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid email or password",
			err,
		)
		return
	}

	match, err := auth.CheckPasswordHash(body.Password, user.HashedPassword)

	if !match {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid email or password",
			err,
		)
		return
	}

	respondWithJSON(w, http.StatusOK, toDomainUser(user))
}
