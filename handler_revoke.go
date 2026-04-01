package main

import (
	"net/http"

	"github.com/cfranklin121/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)

	type ReturnVal struct {
	}

	_, err = cfg.db.RevokeToken(r.Context(), refreshTokenString)
	if err != nil {
		respondWithError(w, 500, "Failed to revoke token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, ReturnVal{})
}
