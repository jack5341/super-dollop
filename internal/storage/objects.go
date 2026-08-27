package storage

import (
	"fmt"
	"io"

	"github.com/minio/minio-go"
)

// ObjectInfo is the subset of object metadata the CLI displays.
type ObjectInfo struct {
	Key  string
	Size int64
}

// Put uploads r under key with the given size (-1 for unknown).
func (c *Client) Put(key string, r io.Reader, size int64) error {
	_, err := c.minio.PutObject(c.cfg.Bucket, key, r, size, minio.PutObjectOptions{
		ContentType: "application/pgp-encrypted",
	})
	if err != nil {
		return fmt.Errorf("upload %s: %w", key, err)
	}
	return nil
}

// Get downloads key as a stream. The caller must close the reader.
func (c *Client) Get(key string) (io.ReadCloser, error) {
	obj, err := c.minio.GetObject(c.cfg.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", key, err)
	}
	return obj, nil
}

// Remove deletes key.
func (c *Client) Remove(key string) error {
	if err := c.minio.RemoveObject(c.cfg.Bucket, key); err != nil {
		return fmt.Errorf("delete %s: %w", key, err)
	}
	return nil
}

// Stat reports whether key exists.
func (c *Client) Stat(key string) error {
	_, err := c.minio.StatObject(c.cfg.Bucket, key, minio.StatObjectOptions{})
	return err
}

// List returns all objects under prefix.
func (c *Client) List(prefix string) ([]ObjectInfo, error) {
	var objects []ObjectInfo
	for obj := range c.minio.ListObjects(c.cfg.Bucket, prefix, true, nil) {
		if obj.Err != nil {
			return nil, fmt.Errorf("list %s: %w", prefix, obj.Err)
		}
		objects = append(objects, ObjectInfo{Key: obj.Key, Size: obj.Size})
	}
	return objects, nil
}
