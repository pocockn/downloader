package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/pocockn/downloader/models"
	"github.com/pocockn/downloader/store"
	"github.com/pocockn/downloader/worker"
)

// Handlers deals with the incoming requests to the API.
type Handlers struct {
	store store.Store
	pool  *worker.Pool
}

// New creates a new Handlers instance to handle requests to the API.
func New(s store.Store, pool *worker.Pool) *Handlers {
	return &Handlers{
		store: s,
		pool:  pool,
	}
}

// URLStore takes a URL and stores it for later processing.
func (h *Handlers) URLStore(c echo.Context) error {
	var url models.URL
	err := c.Bind(&url)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if url.URL == "" {
		return c.String(http.StatusBadRequest, "path must contain url query param")
	}

	h.pool.AddURL(url)

	return nil
}

// URLs returns the latest 50 URLs that the API has received.
func (h *Handlers) URLs(c echo.Context) error {
	bytes, err := h.store.GetAll()
	if err != nil {
		return c.String(http.StatusInternalServerError, "unable fetch urls from the db")
	}

	var urls []models.URL
	for _, d := range bytes {
		var url models.URL
		err := json.Unmarshal(d, &url)
		if err != nil {
			return c.String(http.StatusInternalServerError, "unable to unmarshal bytes into URL")
		}

		urls = append(urls, url)
	}

	if len(urls) <= 50 {
		return c.JSON(http.StatusOK, urls)
	}

	return c.JSON(http.StatusOK, urls[:49])
}
