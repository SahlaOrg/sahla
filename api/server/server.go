package server

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mohamed2394/sahla/api/routes"
	"github.com/mohamed2394/sahla/modules/user/handler"
	"github.com/mohamed2394/sahla/modules/user/repository"
	"github.com/mohamed2394/sahla/pkg/db"
)

type Server struct {
	Echo        *echo.Echo
	UserHandler *handler.UserHandler
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

	userRepository := repository.NewUserRepository(database)
	userHandler := handler.NewUserHandler(userRepository)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register routes
	routes.RegisterUserRoutes(e, userHandler)

	return &Server{
		Echo:        e,
		UserHandler: userHandler,
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
