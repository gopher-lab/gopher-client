package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/masa-finance/gopher-client/types"
)

// LoadConfig loads the Config from environment variables.
// It first attempts to load from a .env file, then falls back to system environment variables.
// The .env file loading is optional - if the file doesn't exist, it will continue without error.
func LoadConfig() (*types.Config, error) {
	// Try to load .env file (optional)
	if err := godotenv.Load(); err != nil {
		// .env file not found or couldn't be loaded, continue with system env vars
		log.Printf("Warning: Could not load .env file: %v", err)
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
		log.Fatalf("Failed to load search configuration: %v", err)
	}
	return config
}
