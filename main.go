package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server running on port", s.Addr)
	log.Fatal(s.ListenAndServe())
}
