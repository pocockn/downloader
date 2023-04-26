package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config struct for config.
type Config struct {
	Port          string        `yaml:"port"`
	Host          string        `yaml:"host"`
	Workers       int           `yaml:"workers"`
	TableName     string        `yaml:"table_name"`
	WatchInterval time.Duration `yaml:"watch_interval"`
}

// New returns a new decoded Config struct
func New(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
