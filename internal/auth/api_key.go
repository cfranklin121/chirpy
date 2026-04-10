package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("%s\n", "No authorization found found")
	}
	keystring := strings.Split(auth, " ")

	if keystring[0] != "ApiKey" {
		return "No Api Key", fmt.Errorf("Error")
	}
	return keystring[1], nil
}
