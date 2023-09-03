package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GET USER API KEY from the HTTP HEADER
// Should look like: Authorization: ApiKey <api_key>
func GetApiKey(h http.Header) (string, error) {

	headerAuth := h.Get("Authorization")

	if headerAuth == "" {
		return "", errors.New("missing AUTH HEADER")
	}

	authValues := strings.Split(headerAuth, " ")

	if len(authValues) != 2 {
		return "", errors.New("malware AUTH HEADER")
	}

	if authValues[0] != "ApiKey" {
		return "", errors.New("malware First part of AUTH HEADER")
	}

	return authValues[1], nil

}
