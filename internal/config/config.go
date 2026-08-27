package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration, sourced from environment
// variables (optionally loaded from a local .env file).
type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	GPGID     string
}

// Load reads configuration from the environment. A missing .env file
// is not an error; missing required variables are.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Endpoint:  os.Getenv("MINIO_ENDPOINT"),
		AccessKey: os.Getenv("MINIO_ACCESS_KEY"),
		SecretKey: os.Getenv("MINIO_SECRET_KEY"),
		Bucket:    os.Getenv("MINIO_BUCKET_NAME"),
		GPGID:     os.Getenv("MINIO_GPG_ID"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var missing []string
	if c.Endpoint == "" {
		missing = append(missing, "MINIO_ENDPOINT")
	}
	if c.AccessKey == "" {
		missing = append(missing, "MINIO_ACCESS_KEY")
	}
	if c.SecretKey == "" {
		missing = append(missing, "MINIO_SECRET_KEY")
	}
	if c.Bucket == "" {
		missing = append(missing, "MINIO_BUCKET_NAME")
	}
	if c.GPGID == "" {
		missing = append(missing, "MINIO_GPG_ID")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}
