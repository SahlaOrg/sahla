package server

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mohamed2394/sahla/api/routes"

	storageHandler "github.com/mohamed2394/sahla/modules/storage/handler"
	"github.com/mohamed2394/sahla/modules/storage/minio"
	storageService "github.com/mohamed2394/sahla/modules/storage/service"

	"github.com/mohamed2394/sahla/modules/user/handler"

	"github.com/mohamed2394/sahla/modules/user/repository"
	"github.com/mohamed2394/sahla/pkg/db"
	"github.com/mohamed2394/sahla/validation"
)

type Server struct {
	Echo           *echo.Echo
	UserHandler    *handler.UserHandler
	StorageService *storageService.StorageService
	StorageHandler *storageHandler.StorageHandler
}

func NewServer(dsn string) (*Server, error) {
	// Connect to the PostgreSQL database
	database, err := db.Connect(dsn)
	if err != nil {
		return nil, err
	}

	// AutoMigrate models
	if err := db.AutoMigrateModels(); err != nil {
		return nil, err
	}

	errV := godotenv.Load()
	if errV != nil {
		log.Fatal("Error loading .env file")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// Initialize MinIO client
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioUseSSL := os.Getenv("MINIO_USE_SSL") == "true"

	minioClient, err := minio.NewMinioClient(minioEndpoint, minioAccessKey, minioSecretKey, minioUseSSL)
	if err != nil {
		return nil, err
	}

	storageService := storageService.NewStorageService(minioClient)

	userRepository := repository.NewUserRepository(database)
	userHandler := handler.NewUserHandler(userRepository)
	storageHandler := storageHandler.NewStorageHandler(storageService, "sahlabucket")

	e := echo.New()
	validation.SetupValidator(e)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes
	routes.RegisterUserRoutes(e, userHandler)
	routes.SetupAuthRoutes(e, userRepository, jwtSecret, refreshSecret)
	routes.RegisterStorageRoutes(e, storageHandler)

	return &Server{
		Echo:           e,
		UserHandler:    userHandler,
		StorageService: storageService,
	}, nil
}

func (s *Server) Start(addr string) {
	log.Println("Server is running at", addr)
	s.Echo.Start(addr)
}

func (s *Server) Close() {
	dbSQL, err := db.GetDB().DB()
	if err != nil {
		log.Fatalf("Error getting db from database: %v", err)
	}
	dbSQL.Close()
}
