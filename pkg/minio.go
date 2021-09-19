package pkg

import (
	"fmt"
	"os"

	"github.com/minio/minio-go"
)

func Connect() *minio.Client {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		fmt.Println(err, "failed while connect to object storage")
	}

	return minioClient
}
