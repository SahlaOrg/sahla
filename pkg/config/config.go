package config

import (
	"os"

	"github.com/mohamed2394/sahla/storage/minio"
)

type AppConfig struct {
	MinioClient *minio.MinioClient
}

func NewAppConfig() (*AppConfig, error) {
	// Initialize MinIO client
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioUseSSL := false 

	minioClient, err := minio.NewMinioClient(minioEndpoint, minioAccessKey, minioSecretKey, minioUseSSL)
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		MinioClient: minioClient,
	}, nil
}
