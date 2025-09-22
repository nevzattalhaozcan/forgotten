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
	var clubRepo repository.ClubRepository = repository.NewClubRepository(s.db)
	var eventRepo repository.EventRepository = repository.NewEventRepository(s.db)
	var bookRepo repository.BookRepository = repository.NewBookRepository(s.db)
	var postRepo repository.PostRepository = repository.NewPostRepository(s.db)

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

	clubService := services.NewClubService(clubRepo, s.config)
	clubHandler := NewClubHandler(clubService)

	eventService := services.NewEventService(eventRepo, clubRepo, s.config)
	eventHandler := NewEventHandler(eventService)

	bookService := services.NewBookService(bookRepo, s.config)
	bookHandler := NewBookHandler(bookService)

	postService := services.NewPostService(postRepo, userRepo, clubRepo, s.config)
	postHandler := NewPostHandler(postService)

	if s.config.Server.Environment != "production" {
		s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		s.router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

	api := s.router.Group("/api/v1")
	{
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)
		api.GET("/clubs", clubHandler.GetAllClubs)
		api.GET("/clubs/:id", clubHandler.GetClub)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(s.config))
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.GET("/users/:id", middleware.AuthorizeSelf(), userHandler.GetUser)
		protected.GET("/users", middleware.RestrictToRoles("admin", "superuser"), userHandler.GetAllUsers)
		protected.PUT("/users/:id", middleware.AuthorizeSelf(), userHandler.UpdateUser)
		protected.DELETE("/users/:id", middleware.AuthorizeSelf(), userHandler.DeleteUser)

		protected.POST("/clubs", clubHandler.CreateClub)
		protected.PUT("/clubs/:id", clubHandler.UpdateClub)
		protected.DELETE("/clubs/:id", clubHandler.DeleteClub)

		protected.POST("/clubs/:id/join", clubHandler.JoinClub)
		protected.POST("/clubs/:id/leave", clubHandler.LeaveClub)
		
		protected.GET("/clubs/:id/members", clubHandler.ListClubMembers)
		protected.PUT("/clubs/:id/members/:user_id", clubHandler.UpdateClubMember)
		protected.GET("/clubs/:id/members/:user_id", clubHandler.GetClubMember)

		protected.POST("/clubs/:id/events", eventHandler.CreateEvent)
		protected.GET("/clubs/:id/events", eventHandler.GetClubEvents)
		protected.GET("/events/:id", eventHandler.GetEvent)
		protected.PUT("/events/:id", eventHandler.UpdateEvent)
		protected.DELETE("/events/:id", eventHandler.DeleteEvent)

		protected.POST("/events/:id/rsvp", eventHandler.RSVPToEvent)
		protected.GET("/events/:id/attendees", eventHandler.GetEventAttendees)

		protected.POST("/books", bookHandler.CreateBook)
		protected.GET("/books/:id", bookHandler.GetBookByID)
		protected.PUT("/books/:id", bookHandler.UpdateBook)
		protected.DELETE("/books/:id", bookHandler.DeleteBook)
		protected.GET("/books", bookHandler.ListBooks)

		protected.POST("/posts", postHandler.CreatePost)
		protected.GET("/posts/:id", postHandler.GetPostByID)
		protected.PUT("/posts/:id", postHandler.UpdatePost)
		protected.DELETE("/posts/:id", postHandler.DeletePost)
		protected.GET("/posts", postHandler.ListAllPosts)

		protected.POST("/posts/:id/like", postHandler.LikePost)
		protected.POST("/posts/:id/unlike", postHandler.UnlikePost)
		protected.GET("/posts/:id/likes", postHandler.ListLikesByPostID)
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}