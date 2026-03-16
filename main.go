package main

import (
	"log"
	"net/http"
)

func main() {
	const filePath = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePath))))
	mux.HandleFunc("/", chirpyHandler)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("Server running on port", port)
	log.Fatal(s.ListenAndServe())
}

func chirpyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
