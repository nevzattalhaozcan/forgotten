package handlers

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/middleware"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
	config *config.Config
	router *gin.Engine
}

func NewServer(db *gorm.DB, config *config.Config) *Server {
	server := &Server{
		db:     db,
		config: config,
		router: gin.New(),
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.Use(gin.Recovery())

	s.router.Use(cors.Default())
	s.router.Use(middleware.LoggingMiddleware())
	s.router.Use(middleware.MetricsMiddleware())

	s.router.HEAD("/health", func(c *gin.Context) {	c.Status(http.StatusOK)	})

	userRepo := repository.NewUserRepository(s.db)
	userService := services.NewUserService(userRepo, s.config)
	userHandler := NewUserHandler(userService)

	api := s.router.Group("/api/v1")
	{
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(s.config))
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.GET("/users/:id", middleware.AuthorizeSelf(), userHandler.GetUser)
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}