package worker_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/pocockn/downloader/mocks"
	"github.com/pocockn/downloader/models"
	"github.com/pocockn/downloader/worker"
)

func TestPool(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)
	urlChan := make(chan models.URL)

	pool := worker.NewPool(3, store, urlChan)

	go pool.Run()
	defer close(urlChan)

	t.Run("New URLs are saved by the workers", func(t *testing.T) {
		url := models.URL{URL: "http://www.example.com"}

		store.EXPECT().Get(url.URL).Return(nil, nil)

		httpmock.RegisterResponder(
			"GET",
			url.URL,
			httpmock.NewStringResponder(200, ``),
		)

		worker.Now = time.Now().UTC()
		url.CreatedAt = worker.Now
		url.Submitted = 1
		bytes, err := json.Marshal(url)
		require.NoError(t, err)
		store.EXPECT().Set(url.URL, bytes).Return(nil)

		pool.AddURL(url)

		time.Sleep(1 * time.Second)
	})

	t.Run("URLs that error are not saved", func(t *testing.T) {
		url := models.URL{URL: "https://www.error.com"}

		httpmock.RegisterResponder(
			"GET",
			url.URL,
			httpmock.NewErrorResponder(fmt.Errorf("big error")),
		)

		pool.AddURL(url)

		time.Sleep(1 * time.Second)
	})

	t.Run("Valid URLs we've seen have their seen number increased", func(t *testing.T) {
		url := models.URL{URL: "https://www.test.com"}
		bytes, err := json.Marshal(url)
		require.NoError(t, err)

		store.EXPECT().Get(url.URL).Return(bytes, nil)

		// Submitted should be +1 from the old value and updated at set to our fixed time.Now
		worker.Now = time.Now().UTC()
		url.Submitted++
		url.UpdatedAt = worker.Now
		updatedBytes, err := json.Marshal(url)
		require.NoError(t, err)
		store.EXPECT().Set(url.URL, updatedBytes).Return(nil)

		httpmock.RegisterResponder(
			"GET",
			url.URL,
			httpmock.NewStringResponder(200, ``),
		)

		pool.AddURL(url)

		time.Sleep(2 * time.Second)
	})
}
