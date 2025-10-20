package config

import (
	"os"
	"path/filepath"

	"github.com/gopher-lab/gopher-client/log"

	"github.com/gopher-lab/gopher-client/types"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// LoadConfig loads the Config from environment variables.
// It first attempts to load from a .env file, then falls back to system environment variables.
// The .env file loading is optional - if the file doesn't exist, it will continue without error.
func LoadConfig() (*types.Config, error) {
	// Try to load .env file (optional)
	// Look for .env file in current directory and parent directories
	envFiles := []string{".env", "../.env", "../../.env"}
	var loaded bool
	for _, envFile := range envFiles {
		if err := godotenv.Load(envFile); err == nil {
			loaded = true
			break
		}
	}

	if !loaded {
		// Try to find .env file by walking up the directory tree
		wd, _ := os.Getwd()
		for {
			envPath := filepath.Join(wd, ".env")
			if _, err := os.Stat(envPath); err == nil {
				if err := godotenv.Load(envPath); err == nil {
					loaded = true
					break
				}
			}
			parent := filepath.Dir(wd)
			if parent == wd {
				break // reached root directory
			}
			wd = parent
		}
	}

	if !loaded {
		log.Warn("Could not find .env file in current or parent directories")
	}

	var config types.Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// MustLoadConfig loads the Config and panics if there's an error.
// Use this when you want the application to fail fast if configuration can't be loaded.
func MustLoadConfig() *types.Config {
	config, err := LoadConfig()
	if err != nil {
		log.Warn("Failed to load search configuration", "error", err)
	}
	return config
}
