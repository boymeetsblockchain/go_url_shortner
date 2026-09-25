// Package api provides the HTTP handlers for the URL shortener.
package api

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/boymeetsblockchain/url_shortner/internal/shortener"
	"github.com/boymeetsblockchain/url_shortner/internal/store"
	"github.com/gin-gonic/gin"
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

func (h *Handler) Routes() *gin.Engine {
	router := gin.Default()
	router.POST("/shorten", h.shorten)
	router.GET("/:code", h.redirect)
	return router
}

func (h *Handler) shorten(c *gin.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
		c.String(http.StatusBadRequest, "invalid request body")
		return
	}

	id := h.counter.Add(1)
	code := shortener.ToBase62(int(id))
	h.store.Set(code, req.URL)

	resp := shortenResponse{
		Code:     code,
		ShortURL: fmt.Sprintf("%s/%s", h.baseURL, code),
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) redirect(c *gin.Context) {
	code := c.Param("code")
	url, ok := h.store.Get(code)
	if !ok {
		c.String(http.StatusNotFound, "short URL not found")
		return
	}
	h.store.IncrementHits(code)
	c.Redirect(http.StatusFound, url)
}
