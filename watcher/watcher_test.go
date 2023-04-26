package watcher_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/pocockn/downloader/mocks"
	"github.com/pocockn/downloader/models"
	"github.com/pocockn/downloader/watcher"
)

func TestWatcher(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mocks.NewMockStore(ctrl)

	w := watcher.New(5*time.Second, store)

	t.Run("Watcher runs every 5 seconds and downloads top 10 submitted URLs", func(t *testing.T) {
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

		for _, url := range urls {
			httpmock.RegisterResponder(
				"GET",
				url.URL,
				httpmock.NewStringResponder(200, ``),
			)
		}

		store.EXPECT().GetAll().Return(results, nil)
		go w.Process()
		time.Sleep(5 * time.Second)
		w.Stop()
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
