package service

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	mClient "github.com/mohamed2394/sahla/storage/minio"
)

type StorageService struct {
	minioClient *mClient.MinioClient
}

type FileInfo struct {
	Name         string
	Size         int64
	LastModified time.Time
	ContentType  string
}

func NewStorageService(minioClient *mClient.MinioClient) *StorageService {
	return &StorageService{minioClient: minioClient}
}

func (s *StorageService) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) error {
	return s.minioClient.UploadFile(ctx, bucketName, objectName, reader, size, contentType)
}

func (s *StorageService) DownloadFile(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	return s.minioClient.DownloadFile(ctx, bucketName, objectName)
}

func (s *StorageService) ListFiles(ctx context.Context, bucketName string) ([]FileInfo, error) {
	objects := s.minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})
	var files []FileInfo

	for object := range objects {
		if object.Err != nil {
			return nil, object.Err
		}
		files = append(files, FileInfo{
			Name:         object.Key,
			Size:         object.Size,
			LastModified: object.LastModified,
			ContentType:  object.ContentType,
		})
	}

	return files, nil
}

func (s *StorageService) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	return s.minioClient.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}

func (s *StorageService) GetFileInfo(ctx context.Context, bucketName, objectName string) (*FileInfo, error) {
	info, err := s.minioClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:         info.Key,
		Size:         info.Size,
		LastModified: info.LastModified,
		ContentType:  info.ContentType,
	}, nil
}

func (s *StorageService) FileExists(ctx context.Context, bucketName, objectName string) (bool, error) {
	_, err := s.minioClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
