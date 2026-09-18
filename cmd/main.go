package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

var store = NewStore()
var counter atomic.Int64 // thread-safe incrementing ID

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id := counter.Add(1)
	code := toBase62(int(id))
	store.Set(code, req.URL)

	resp := shortenResponse{
		Code:     code,
		ShortURL: fmt.Sprintf("http://localhost:8080/%s", code),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code") // Go 1.22+ wildcard routing
	url, ok := store.Get(code)
	if !ok {
		http.Error(w, "short URL not found", http.StatusNotFound)
		return
	}
	store.IncrementHits(code)
	http.Redirect(w, r, url, http.StatusFound) // 302
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", shortenHandler)
	mux.HandleFunc("GET /{code}", redirectHandler)

	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", mux)
}
