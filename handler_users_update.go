package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cfranklin121/chirpy/internal/auth"
	"github.com/cfranklin121/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	defer r.Body.Close()
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

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Token")
		return
	}

	hashedPassword, err := auth.HashPassword(reqBody.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password")
		return
	}

	params := database.UpdateUserParams{
		ID:             userID,
		Email:          reqBody.Email,
		HashedPassword: hashedPassword,
	}

	updatedUser, err := cfg.db.UpdateUser(r.Context(), params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, ReturnVal{
		User: User{
			ID:        updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email:     updatedUser.Email,
		},
	})

}
