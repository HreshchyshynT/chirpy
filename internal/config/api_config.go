package config

import (
	"net/http"
	"sync/atomic"

	"github.com/hreshchyshynt/chirpy/internal/database"
)

type DbHandler func(
	w http.ResponseWriter,
	r *http.Request,
	queries *database.Queries,
)

type ApiHandler func(
	w http.ResponseWriter,
	r *http.Request,
	apiConfig *ApiConfig,
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Queries        *database.Queries
	Platform       Platform
	Secret         string
}
