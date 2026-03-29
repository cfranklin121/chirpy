package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cfranklin121/chirpy/internal/auth"
	"github.com/cfranklin121/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("%s %s", r.Method, r.URL.Path)
	type RequestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type ReturnVal struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	reqBody := RequestBody{}
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondWithError(w, 500, "Could not decode")
		return
	}
	expirationTime := time.Hour
	log.Println(reqBody)

	expirationTime = time.Duration(60) * time.Second

	user, err := cfg.db.GetUser(r.Context(), reqBody.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(reqBody.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	if !match {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(expirationTime))
	if err != nil {
		respondWithError(w, 500, "Could not create access token")
		return
	}

	refreshParams := database.CreateRefreshTokenParams{
		Token:     auth.MakeRefreshToken(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(1440) * time.Hour),
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		respondWithError(w, 500, "Could not create refresh token")
		return
	}

	respondWithJSON(w, 200, ReturnVal{
		User: User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			AccessToken:  accessToken,
			RefreshToken: refreshToken.Token,
		},
	})

}
