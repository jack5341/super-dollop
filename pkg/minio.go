package pkg

import (
	"fmt"
	"net/http"
	"os"

	"github.com/minio/minio-go"
)

func Connect() *minio.Client {

	// For checking connection on minIO server can use http://<endpoint>/minio/health/live
	url := fmt.Sprintf("http://%s/minio/health/live", os.Getenv("MINIO_ENDPOINT"))
	_, err := http.Get(url)

	if err != nil {
		panic("failed to connect to minio client")
	}

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
