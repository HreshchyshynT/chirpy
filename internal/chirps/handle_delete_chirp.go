package chirps

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/config"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func HandleDeleteChirp(
	w http.ResponseWriter,
	r *http.Request,
	userId uuid.UUID,
	apiConfig *config.ApiConfig,
) {
	type ResponseBody struct {
	}

	chirpIdStr := r.PathValue("chirpID")
	if len(chirpIdStr) == 0 {
		utils.RespondWithError(
			w,
			http.StatusBadRequest,
			"chirpID is empty",
			errors.New("chirpID is empty"),
		)
		return
	}

	chirpId, err := uuid.Parse(chirpIdStr)

	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusBadRequest,
			"chirpID must be valid UUID",
			err,
		)
		return
	}

	chirp, err := apiConfig.Queries.GetChirp(r.Context(), chirpId)

	// TODO: handle not found
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			utils.MessageInternalServerError,
			err,
		)
		return
	}

	if chirp.UserID != userId {
		utils.RespondWithError(
			w,
			http.StatusForbidden,
			utils.MessageForbidden,
			err,
		)
		return
	}

	err = apiConfig.Queries.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			utils.MessageInternalServerError,
			err,
		)
		return
	}

	utils.RespondWithJSON(
		w,
		http.StatusNoContent,
		ResponseBody{},
	)
}
