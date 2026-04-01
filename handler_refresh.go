package main

import (
	"net/http"
	"time"

	"github.com/cfranklin121/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)

	type ReturnVal struct {
		AccessToken string `json:"token"`
	}

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshTokenString)

	refresh_token, err := cfg.db.GetRefreshToken(r.Context(), userID.ID)
	if err != nil {
		respondWithError(w, 500, "Could not get refresh token")
		return
	}

	if refresh_token.RevokedAt.Valid {
		respondWithError(w, 401, "Token has been revoked")
		return
	}

	if refresh_token.RevokedAt.Valid {
		respondWithError(w, 401, "Token has been revoked")
		return
	}

	expirationTime := time.Hour
	expirationTime = time.Duration(60) * time.Second

	accessToken, err := auth.MakeJWT(userID.ID, cfg.secret, time.Duration(expirationTime))
	if err != nil {
		respondWithError(w, 500, "Could not create access token")
		return
	}

	respondWithJSON(w, 200, ReturnVal{
		AccessToken: accessToken,
	})
}
