package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		respondWithError(w, 500, "Invalid authorization")
		return
	}

	type ReturnVal struct {
	}

	tknstrng := strings.Split(authorization, " ")

	_, err := cfg.db.RevokeToken(r.Context(), tknstrng[1])
	if err != nil {
		respondWithError(w, 500, "Failed to revoke token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, ReturnVal{})
}
