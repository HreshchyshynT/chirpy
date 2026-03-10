package chirps

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func HandleGetChirp(
	w http.ResponseWriter,
	r *http.Request,
	db *database.Queries,
) {
	idStr := r.PathValue("chirpID")

	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.RespondWithError(w,
			http.StatusInternalServerError,
			"Can't parse id",
			err,
		)
		return
	}
	chirp, err := db.GetChirp(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w,
			http.StatusNotFound,
			"Chirp not found",
			err,
		)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, toDomainChirp(chirp))
}
