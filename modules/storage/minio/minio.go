package minio

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
}

func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &MinioClient{client: client}, nil
}

func (m *MinioClient) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{ContentType: contentType})
	return err
}

func (m *MinioClient) DownloadFile(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	object, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	return io.ReadAll(object)
}

func (m *MinioClient) ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo {
	return m.client.ListObjects(ctx, bucketName, opts)
}

func (m *MinioClient) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, opts)
}

func (m *MinioClient) StatObject(ctx context.Context, bucketName, objectName string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	return m.client.StatObject(ctx, bucketName, objectName, opts)
}
