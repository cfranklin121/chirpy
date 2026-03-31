package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/cfranklin121/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		respondWithError(w, 500, "Invalid authorization")
		return
	}

	type ReturnVal struct {
		Token string
	}

	tknstrng := strings.Split(authorization, " ")

	refresh_token, err := cfg.db.GetRefreshToken(r.Context(), tknstrng[1])
	if err != nil {
		respondWithError(w, 500, "Could not get refresh token")
		return
	}

	if refresh_token.RevokedAt.Valid {
		respondWithError(w, 401, "Token has been revoked")
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), refresh_token.Token)

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
		Token: accessToken,
	})
}
