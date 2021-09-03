package pkg

import (
	"github.com/minio/minio-go"
)

func Connect() *minio.Client {
	endpoint := "superdollop"
	accessKeyID := "s3accesskey"
	secretAccessKey := "s3secretkey"
	useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		panic(err)
	}

	return minioClient
}
