package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	keyValue := headers.Get("Authorization")
	if len(keyValue) == 0 || !strings.HasPrefix(keyValue, "ApiKey") {
		return "", errors.New("No api key provided")
	}

	splitted := strings.Split(keyValue, " ")
	return splitted[1], nil
}
