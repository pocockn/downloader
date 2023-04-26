package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pocockn/downloader/config"
)

func TestNewConfig(t *testing.T) {
	cfg, err := config.New("../config.yaml")
	require.NoError(t, err)
	assert.Equal(t, 3, cfg.Workers)
	assert.Equal(t, "5000", cfg.Port)
	assert.Equal(t, "urls", cfg.TableName)
	assert.Equal(t, "127.0.0.1", cfg.Host)
	assert.Equal(t, 60*time.Second, cfg.WatchInterval)
}
