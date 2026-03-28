package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cfranklin121/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	log.Printf("%s %s", r.Method, r.URL.Path)
	type RequestBody struct {
		Password  string `json:"password"`
		Email     string `json:"email"`
		ExpiresIn int    `json:"expires_in_seconds"`
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
	log.Println(reqBody)
	if reqBody.ExpiresIn == 0 {
		reqBody.ExpiresIn = 60
	}

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

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(reqBody.ExpiresIn))

	respondWithJSON(w, 200, ReturnVal{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		},
	})

}
