package chirps

import (
	"net/http"

	"github.com/hreshchyshynt/chirpy/internal/database"
	"github.com/hreshchyshynt/chirpy/internal/utils"
)

func HandleGetAllChirps(
	w http.ResponseWriter,
	r *http.Request,
	db *database.Queries,
) {

	chirps, err := db.GetAllChirps(r.Context())
	if err != nil {
		utils.RespondWithError(
			w,
			http.StatusInternalServerError,
			"Can't get chirps",
			err,
		)
		return
	}

	responseChirps := make([]Chirp, len(chirps))

	for i, c := range chirps {
		responseChirps[i] = toDomainChirp(c)
	}

	utils.RespondWithJSON(w, http.StatusOK, responseChirps)
}
