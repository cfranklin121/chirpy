package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		return
	}

	type ReturnVal struct {
		Token string
	}

	//tknstrng := strings.Split(authorization, " ")
}
