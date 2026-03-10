package user

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/auth"
	"github.com/hreshchyshynt/chirpy/internal/config"
	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func HandleUpdateUser(
	w http.ResponseWriter,
	r *http.Request,
	userId uuid.UUID,
	apiConfig *config.ApiConfig,
) {
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	decoder.DisallowUnknownFields()

	request := requestBody{}

	err := decoder.Decode(&request)

	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusBadRequest,
			utils.MessageInvalidRequestBody,
			err,
		)
		return
	}

	hashedPassword, err := auth.HashPassword(request.Password)

	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			utils.MessageInternalServerError,
			err,
		)
		return
	}

	user, err := apiConfig.Queries.UpdateUser(r.Context(), database.UpdateUserParams{
		UserID:         userId,
		Email:          request.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			utils.MessageInternalServerError,
			err,
		)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, toDomainUser(user))
}
