package config

import (
	"fmt"
	"os"
	"time"

	yaml "github.com/goccy/go-yaml"
)

func New(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	replaced := os.ExpandEnv(string(data))

	var config Config
	err = yaml.Unmarshal([]byte(replaced), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml config: %w", err)
	}

	return &config, nil
}

type Config struct {
	Port     int       `yaml:"port"`
	Postgres *Postgres `yaml:"postgres"`
}

type Postgres struct {
	Timeout  time.Duration `yaml:"timeout"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Database string        `yaml:"database"`
	URL      string        `yaml:"url"`
}
