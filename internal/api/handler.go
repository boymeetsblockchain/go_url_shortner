// Package api provides the HTTP handlers for the URL shortener.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/boymeetsblockchain/url_shortner/internal/shortener"
	"github.com/boymeetsblockchain/url_shortner/internal/store"
)

type Handler struct {
	store   *store.Store
	counter atomic.Int64
	baseURL string
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

func NewHandler(s *store.Store, baseURL string) *Handler {
	return &Handler{store: s, baseURL: baseURL}
}

func (h *Handler) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.shorten)
	mux.HandleFunc("GET /{code}", h.redirect)
	return mux
}

func (h *Handler) shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id := h.counter.Add(1)
	code := shortener.ToBase62(int(id))
	h.store.Set(code, req.URL)

	resp := shortenResponse{
		Code:     code,
		ShortURL: fmt.Sprintf("%s/%s", h.baseURL, code),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code") // Go 1.22+ wildcard routing
	url, ok := h.store.Get(code)
	if !ok {
		http.Error(w, "short URL not found", http.StatusNotFound)
		return
	}
	h.store.IncrementHits(code)
	http.Redirect(w, r, url, http.StatusFound) // 302
}
