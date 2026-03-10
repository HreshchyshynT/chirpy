package utils

import (
	"errors"

	"github.com/lib/pq"
)

func IsDuplicatedKeys(err error) bool {
	var pqErr *pq.Error
	// duplicated keys error has code 23505
	return err != nil && errors.As(err, &pqErr) && pqErr.Code == "23505"
}
