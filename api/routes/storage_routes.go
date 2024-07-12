package routes

import (
	"github.com/labstack/echo/v4"
	handler "github.com/mohamed2394/sahla/storage/handler"

)

func RegisterStorageRoutes(e *echo.Echo, h *handler.StorageHandler) {
	e.POST("/upload", h.UploadFile)
	e.GET("/download/:filename", h.DownloadFile)
	e.GET("/files", h.ListFiles)
	e.DELETE("/files/:filename", h.DeleteFile)
	e.GET("/files/:filename/info", h.GetFileInfo)
	e.GET("/files/:filename/exists", h.FileExists)
}
