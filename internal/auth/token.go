package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	value := headers.Get("Authorization")
	if len(value) == 0 {
		return "", errors.New("No Authorization header provided")
	}

	if !strings.HasPrefix(value, "Bearer ") {
		return "", errors.New("No Bearer token provided")
	}

	splitted := strings.Split(value, " ")
	if len(splitted) != 2 {
		return "", errors.New("No valid token provided")
	}

	return splitted[1], nil
}
