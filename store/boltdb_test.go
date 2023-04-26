package store_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pocockn/downloader/models"
	"github.com/pocockn/downloader/store"
)

func TestBolt(t *testing.T) {
	var db store.Store
	var err error

	t.Run("Can connect to Bolt", func(t *testing.T) {
		db, err = store.ConnectBolt("test")
		require.NoError(t, err)
	})

	t.Run("Set item within database", func(t *testing.T) {
		assert.NoError(t, db.Set("test", []byte("test_bytes")))
	})

	t.Run("Get item within database", func(t *testing.T) {
		result, err := db.Get("test")
		require.NoError(t, err)

		assert.Equal(t, []byte("test_bytes"), result)
	})

	assert.NoError(t, db.Disconnect())
	assert.NoError(t, os.Remove("my.db"))
}

func TestBoltDB_GetAll(t *testing.T) {
	db, err := store.ConnectBolt("test")
	require.NoError(t, err)

	for i := 0; i <= 10; i++ {
		url := models.URL{URL: fmt.Sprintf("www.example.com%d", i)}
		urlBytes, err := json.Marshal(url)
		assert.NoError(t, err)
		assert.NoError(t, db.Set(url.URL, urlBytes))
	}

	result, err := db.GetAll()
	assert.NoError(t, err)
	assert.Len(t, result, 11)

	assert.NoError(t, db.Disconnect())
	assert.NoError(t, os.Remove("my.db"))
}
