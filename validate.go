package main

import (
	"errors"
	"strings"
)

func validateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", errors.New("Chirp is too long")
	}

	cleanedBody := cleanText(chirp)
	return cleanedBody, nil
}

func cleanText(input string) string {
	var builder strings.Builder

	profanes := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	splitted := strings.Split(input, " ")

	for i, s := range splitted {

		if _, ok := profanes[strings.ToLower(s)]; ok {
			builder.WriteString("****")
		} else {
			builder.WriteString(s)
		}

		if i < len(splitted)-1 {
			builder.WriteRune(' ')
		}
	}

	return builder.String()
}
