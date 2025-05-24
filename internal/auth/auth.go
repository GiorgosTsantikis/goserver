package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Authorization: ApiKey {insert api key}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("missing Authorization header")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 || vals[0] != "ApiKey" || len(vals[1]) != 64 {
		return "", errors.New("malformed Authorization header")
	}
	return vals[1], nil
}
