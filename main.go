package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	w.Header().Add("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
  </body>
</html>
		`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	log.Printf("%s %s", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	s := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Println("Server running on port", port)
	log.Fatal(s.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
	w.Header().Add("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

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
