package handlers

import (
	"net/http"
	"time"

	_ "github.com/nevzattalhaozcan/forgotten/docs" // swag doc import
	"github.com/nevzattalhaozcan/forgotten/pkg/cache"
	"github.com/nevzattalhaozcan/forgotten/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

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

	var userRepo repository.UserRepository = repository.NewUserRepository(s.db)

	var rdbAvailable bool
	var ttl time.Duration
	var rdb *redis.Client

	if s.config.Redis.Enabled {
		client, err := cache.NewRedisClient(&s.config.Redis)
		if err != nil {
			logger.Warn("Redis enabled but connection failed, continuing without cache")
		} else {
			rdb = client
			rdbAvailable = true
			ttl = time.Duration(s.config.Redis.CacheTTLSeconds) * time.Second
		}
	}

	if rdbAvailable {
		userRepo = repository.NewCachedUserRepository(userRepo, rdb, ttl)
		logger.Info("User repository caching enabled")
	}

	userService := services.NewUserService(userRepo, s.config)
	userHandler := NewUserHandler(userService)

	if s.config.Server.Environment != "production" {
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

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
		protected.GET("/users", middleware.RestrictToRoles("admin", "superuser"), userHandler.GetAllUsers)
		protected.PUT("/users/:id", middleware.AuthorizeSelf(), userHandler.UpdateUser)
		protected.DELETE("/users/:id", middleware.AuthorizeSelf(), userHandler.DeleteUser)
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}