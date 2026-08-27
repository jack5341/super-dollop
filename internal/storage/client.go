package storage

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jack5341/super-dollop/internal/config"
	"github.com/minio/minio-go"
)

const healthTimeout = 5 * time.Second

// Client wraps the MinIO client together with the resolved config.
type Client struct {
	minio *minio.Client
	cfg   *config.Config
}

// NewClient verifies the MinIO server is reachable and healthy, then
// returns a ready-to-use client.
func NewClient(cfg *config.Config) (*Client, error) {
	healthURL := fmt.Sprintf("http://%s/minio/health/live", cfg.Endpoint)
	httpClient := &http.Client{Timeout: healthTimeout}

	resp, err := httpClient.Get(healthURL)
	if err != nil {
		return nil, fmt.Errorf("cannot reach minio at %s: %w", cfg.Endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("minio health check at %s returned status %d", cfg.Endpoint, resp.StatusCode)
	}

	minioClient, err := minio.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey, false)
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	return &Client{minio: minioClient, cfg: cfg}, nil
}
