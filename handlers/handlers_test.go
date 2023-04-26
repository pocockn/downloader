package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pocockn/downloader/handlers"
	"github.com/pocockn/downloader/mocks"
	"github.com/pocockn/downloader/models"
	"github.com/pocockn/downloader/worker"
)

func TestURLStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	urlsChan := make(chan models.URL, 10)

	h := handlers.New(store, worker.NewPool(3, store, urlsChan))

	t.Run("store endpoint must contain url query param", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/store", http.NoBody)
		assert.NoError(t, err)

		e := echo.New()
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = h.URLStore(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("URL is passed to a worker", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/store?url=http://www.example.com", http.NoBody)
		assert.NoError(t, err)

		e := echo.New()
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = h.URLStore(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestURLs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	urlsChan := make(chan models.URL, 10)

	h := handlers.New(store, worker.NewPool(3, store, urlsChan))

	t.Run("URLs endpoint returns up to 50 of the latest URLs", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/urls", http.NoBody)
		assert.NoError(t, err)

		urls := []models.URL{
			{URL: "http://www.example.com", Submitted: 0},
			{URL: "http://www.example1.com", Submitted: 1},
			{URL: "http://www.example2.com", Submitted: 2},
			{URL: "http://www.example3.com", Submitted: 3},
			{URL: "http://www.example4.com", Submitted: 4},
			{URL: "http://www.example5.com", Submitted: 5},
			{URL: "http://www.example6.com", Submitted: 6},
			{URL: "http://www.example7.com", Submitted: 7},
			{URL: "http://www.example8.com", Submitted: 8},
			{URL: "http://www.example9.com", Submitted: 9},
			{URL: "http://www.example10.com", Submitted: 10},
		}

		results := marshalURLs(urls, t)

		store.EXPECT().GetAll().Return(results, nil)

		e := echo.New()
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err = h.URLs(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func marshalURLs(urls []models.URL, t *testing.T) [][]byte {
	var results [][]byte
	for _, url := range urls {
		bytes, err := json.Marshal(url)
		require.NoError(t, err)
		results = append(results, bytes)
	}

	return results
}
