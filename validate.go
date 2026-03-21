package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type RequestBody struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	dat, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("%s %s", r.Method, r.URL.Path)
		respondWithError(w, 500, "Could not read data")
		return
	}

	reqBody := RequestBody{}
	err = json.Unmarshal(dat, &reqBody)
	if err != nil {
		log.Printf("%s %s", r.Method, r.URL.Path)
		respondWithError(w, 500, "Could not unmarshal JSON")
		return
	}

	if len(reqBody.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	log.Printf("%s %s", r.Method, r.URL.Path)
	cleaned := cleanString(reqBody.Body)
	respondWithJSON(w, 200, returnVal{
		Cleaned_body: cleaned,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func cleanString(str string) string {
	slice := strings.Split(str, " ")
	newStr := []string{}
	for _, word := range slice {
		switch strings.ToLower(word) {
		case "kerfuffle":
			word = "****"
		case "sharbert":
			word = "****"
		case "fornax":
			word = "****"
		}
		newStr = append(newStr, word)
	}
	return strings.Join(newStr, " ")
}
