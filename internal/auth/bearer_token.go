package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("%s\n", "No authorization found found")
	}
	tknstrng := strings.Split(auth, " ")

	if tknstrng[0] != "Bearer" {
		return "No bearer token", fmt.Errorf("Error")
	}

	return tknstrng[1], nil
}
