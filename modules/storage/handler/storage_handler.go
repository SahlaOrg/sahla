package handler

import (
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/mohamed2394/sahla/modules/storage/service"
)

type StorageHandler struct {
	storageService *service.StorageService
	bucketName     string
}

func NewStorageHandler(storageService *service.StorageService, bucketName string) *StorageHandler {
	return &StorageHandler{
		storageService: storageService,
		bucketName:     bucketName,
	}
}

func (h *StorageHandler) UploadFile(c echo.Context) error {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to get file from request"})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open file"})
	}
	defer src.Close()

	// Get the file size
	size := file.Size

	// Get the file name and extension
	filename := filepath.Base(file.Filename)

	// Determine content type
	contentType := file.Header.Get("Content-Type")

	// Upload the file
	err = h.storageService.UploadFile(c.Request().Context(), h.bucketName, filename, src, size, contentType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to upload file"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "File uploaded successfully", "filename": filename})
}

func (h *StorageHandler) DownloadFile(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filename is required"})
	}

	data, err := h.storageService.DownloadFile(c.Request().Context(), h.bucketName, filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to download file"})
	}

	// Determine the content type
	contentType := http.DetectContentType(data)

	return c.Blob(http.StatusOK, contentType, data)
}

func (h *StorageHandler) ListFiles(c echo.Context) error {
	files, err := h.storageService.ListFiles(c.Request().Context(), h.bucketName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to list files"})
	}

	return c.JSON(http.StatusOK, files)
}

func (h *StorageHandler) DeleteFile(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filename is required"})
	}

	err := h.storageService.DeleteFile(c.Request().Context(), h.bucketName, filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete file"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "File deleted successfully"})
}

func (h *StorageHandler) GetFileInfo(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filename is required"})
	}

	fileInfo, err := h.storageService.GetFileInfo(c.Request().Context(), h.bucketName, filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get file info"})
	}

	return c.JSON(http.StatusOK, fileInfo)
}

func (h *StorageHandler) FileExists(c echo.Context) error {
	filename := c.Param("filename")
	if filename == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Filename is required"})
	}

	exists, err := h.storageService.FileExists(c.Request().Context(), h.bucketName, filename)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to check file existence"})
	}

	return c.JSON(http.StatusOK, map[string]bool{"exists": exists})
}
