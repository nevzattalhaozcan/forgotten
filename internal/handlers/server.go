package handlers

import (
	"net/http"
	"time"

	_ "github.com/nevzattalhaozcan/forgotten/docs" // swag doc import
	"github.com/nevzattalhaozcan/forgotten/pkg/cache"
	"github.com/nevzattalhaozcan/forgotten/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nevzattalhaozcan/forgotten/internal/clients"
	"github.com/nevzattalhaozcan/forgotten/internal/config"
	"github.com/nevzattalhaozcan/forgotten/internal/middleware"
	"github.com/nevzattalhaozcan/forgotten/internal/repository"
	"github.com/nevzattalhaozcan/forgotten/internal/services"
	"gorm.io/gorm"
)

type Server struct {
	db     *gorm.DB
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

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://forgotten.onrender.com", "https://forgotten-ui.vercel.app", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	s.router.Use(middleware.LoggingMiddleware())
	s.router.Use(middleware.MetricsMiddleware())
	_ = s.router.SetTrustedProxies(nil)

	s.router.HEAD("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	var userRepo repository.UserRepository = repository.NewUserRepository(s.db)
	var clubRepo repository.ClubRepository = repository.NewClubRepository(s.db)
	var eventRepo repository.EventRepository = repository.NewEventRepository(s.db)
	var bookRepo repository.BookRepository = repository.NewBookRepository(s.db)
	var bookClient clients.BookAPIClient = clients.NewMultiSourceClient(s.config.BookAPIs.GoogleBooksAPIKey, s.config.BookAPIs.PreferredSource)
	var postRepo repository.PostRepository = repository.NewPostRepository(s.db)
	var commentRepo repository.CommentRepository = repository.NewCommentRepository(s.db)
	var readingRepo repository.ReadingRepository = repository.NewReadingRepository(s.db)
	var clubReadingRepo repository.ClubReadingRepository = repository.NewClubReadingRepository(s.db)
	var clubRatingRepo repository.ClubRatingRepository = repository.NewClubRatingRepository(s.db)

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

	clubService := services.NewClubService(clubRepo, clubRatingRepo, s.config)
	clubHandler := NewClubHandler(clubService)

	eventService := services.NewEventService(eventRepo, clubRepo, s.config)
	eventHandler := NewEventHandler(eventService)

	bookService := services.NewBookService(bookRepo, bookClient, s.config)
	bookHandler := NewBookHandler(bookService)

	postService := services.NewPostService(postRepo, userRepo, clubRepo, bookRepo, s.db, s.config)
	postHandler := NewPostHandler(postService)

	commentService := services.NewCommentService(commentRepo, postRepo, userRepo, s.config)
	commentHandler := NewCommentHandler(commentService)

	readingService := services.NewReadingService(s.config, userRepo, bookRepo, clubRepo, readingRepo, clubReadingRepo)
	readingHandler := NewReadingHandler(readingService)

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
		api.GET("/clubs/:id/members", clubHandler.ListClubMembers)
		api.GET("/clubs/:id/ratings", clubHandler.ListClubRatings)

		api.GET("/posts/public", postHandler.ListPublicPosts)
		api.GET("/posts/popular", postHandler.ListPopularPublicPosts)

		api.GET("/books", bookHandler.Search)
		api.GET("/books/:id", bookHandler.GetBookByID)

		api.GET("/posts/:id/likes", postHandler.ListLikesByPostID)
		api.GET("/posts", postHandler.ListAllPosts)
		api.GET("/clubs/:id/posts/summaries", middleware.OptionalAuthMiddleware(s.config), postHandler.ListPostSummaries)
		api.GET("/posts/:id", postHandler.GetPostByID)

		api.GET("/posts/:id/comments", commentHandler.ListCommentsByPostID)
		api.GET("/users/:id/comments", commentHandler.ListCommentsByUserID)
		api.GET("/comments/:id", commentHandler.GetCommentByID)
		api.GET("/comments/:id/likes", commentHandler.ListLikesByCommentID)

		api.GET("/events/public", eventHandler.GetPublicEvents)
	}

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(s.config))
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.GET("/users/:id", middleware.AuthorizeSelf(), userHandler.GetUser)
		protected.GET("/users", middleware.RestrictToRoles("admin", "superuser"), userHandler.GetAllUsers)
		protected.PUT("/users/:id", middleware.AuthorizeSelf(), userHandler.UpdateUser)
		protected.DELETE("/users/:id", middleware.AuthorizeSelf(), userHandler.DeleteUser)
		protected.PATCH("/users/:id/password", userHandler.PatchPassword)
		protected.PATCH("/users/:id/profile", userHandler.PatchProfile)
		protected.PATCH("/users/:id/account", userHandler.PatchAccount)
		protected.PATCH("/users/:id/avatar", userHandler.PatchAvatar)

		protected.POST("/clubs", clubHandler.CreateClub)
		protected.PUT("/clubs/:id", middleware.RequireClubMembershipWithRoles(clubRepo, "club_admin"), clubHandler.UpdateClub)
		protected.DELETE("/clubs/:id", middleware.RequireClubMembershipWithRoles(clubRepo, "club_admin"), clubHandler.DeleteClub)
		protected.GET("/clubs/:id/poll", postHandler.GetPollPostsByClubID)

		protected.POST("/clubs/:id/join", clubHandler.JoinClub)
		protected.POST("/clubs/:id/leave", middleware.RequireClubMembership(clubRepo), clubHandler.LeaveClub)
		protected.POST("/clubs/:id/ratings", middleware.RequireClubMembership(clubRepo), clubHandler.RateClub)
		protected.GET("/my-clubs", clubHandler.GetMyClubs)

		protected.PUT("/clubs/:id/members/:user_id", middleware.RequireClubMembershipWithRoles(clubRepo, "club_admin", "moderator"), clubHandler.UpdateClubMember)
		protected.GET("/clubs/:id/members/:user_id", clubHandler.GetClubMember)

		protected.POST("/clubs/:id/events", middleware.RequireClubMembershipWithRoles(clubRepo, "club_admin", "moderator"), eventHandler.CreateEvent)
		protected.GET("/clubs/:id/events", middleware.RequireClubMembership(clubRepo), eventHandler.GetClubEvents)
		protected.GET("/events/:id", middleware.RequireClubMembership(clubRepo), eventHandler.GetEvent)
		protected.PUT("/events/:id", middleware.RequireClubMembershipWithRoles(clubRepo, "club_admin", "moderator"), eventHandler.UpdateEvent)
		protected.DELETE("/events/:id", middleware.RequireClubMembershipWithRoles(clubRepo, "club_admin", "moderator"), eventHandler.DeleteEvent)

		protected.POST("/events/:id/rsvp", middleware.RequireClubMembership(clubRepo), eventHandler.RSVPToEvent)
		protected.GET("/events/:id/attendees", middleware.RequireClubMembership(clubRepo), eventHandler.GetEventAttendees)

		protected.POST("/books", middleware.RestrictToRoles("admin", "superuser"), bookHandler.CreateBook)
		protected.PUT("/books/:id", middleware.RestrictToRoles("admin", "superuser"), bookHandler.UpdateBook)
		protected.DELETE("/books/:id", middleware.RestrictToRoles("admin", "superuser"), bookHandler.DeleteBook)

		protected.POST("/posts", middleware.RequireClubMembership(clubRepo), postHandler.CreatePost)
		protected.PUT("/posts/:id", middleware.RequireClubMembership(clubRepo), postHandler.UpdatePost)
		protected.DELETE("/posts/:id", middleware.RequireClubMembership(clubRepo), postHandler.DeletePost)
		protected.GET("/posts/reviews", middleware.RequireClubMembership(clubRepo), postHandler.GetReviewsByBook)
		protected.GET("/posts/filter", middleware.RequireClubMembership(clubRepo), postHandler.GetPostsByType)

		protected.POST("/posts/:id/vote", middleware.RequireClubMembership(clubRepo), postHandler.VoteOnPoll)
		protected.POST("/posts/:id/unvote", middleware.RequireClubMembership(clubRepo), postHandler.RemoveVoteFromPoll)
		protected.GET("/posts/:id/poll/votes", middleware.RequireClubMembership(clubRepo), postHandler.GetUserPollVotes)

		protected.POST("/posts/:id/like", middleware.RequireClubMembership(clubRepo), postHandler.LikePost)
		protected.POST("/posts/:id/unlike", middleware.RequireClubMembership(clubRepo), postHandler.UnlikePost)

		protected.POST("/posts/:id/comments", middleware.RequireClubMembership(clubRepo), commentHandler.CreateComment)
		protected.PUT("/comments/:id", middleware.RequireClubMembership(clubRepo), commentHandler.UpdateComment)
		protected.DELETE("/comments/:id", middleware.RequireClubMembership(clubRepo), commentHandler.DeleteComment)

		protected.POST("/comments/:id/like", middleware.RequireClubMembership(clubRepo), commentHandler.LikeComment)
		protected.POST("/comments/:id/unlike", middleware.RequireClubMembership(clubRepo), commentHandler.UnlikeComment)

		protected.POST("/users/:id/reading/sync", middleware.AuthorizeSelf(), readingHandler.SyncUserStats)
		protected.POST("/users/:id/reading/start", middleware.AuthorizeSelf(), readingHandler.StartReading)
		protected.PATCH("/users/:id/reading/:bookID/progress", middleware.AuthorizeSelf(), readingHandler.UpdateProgress)
		protected.POST("/users/:id/reading/:bookID/complete", middleware.AuthorizeSelf(), readingHandler.CompleteReading)
		protected.GET("/users/:id/readings/:bookID/progress", middleware.AuthorizeSelf(), readingHandler.GetProgress)
		protected.GET("/users/:id/reading", middleware.AuthorizeSelf(), readingHandler.ListUserProgress)
		protected.GET("/users/:id/reading/history", readingHandler.UserReadingHistory)

		protected.POST("/clubs/:id/reading/assign", middleware.RequireClubMembershipWithRoles(clubRepo, "admin", "moderator"), readingHandler.AssignBookToClub)
		protected.PATCH("/clubs/:id/reading/checkpoint", middleware.RequireClubMembershipWithRoles(clubRepo, "admin", "moderator"), readingHandler.UpdateClubCheckpoint)
		protected.POST("/clubs/:id/reading/complete", middleware.RequireClubMembershipWithRoles(clubRepo, "admin", "moderator"), readingHandler.CompleteClubAssignment)
		protected.GET("/clubs/:id/reading", readingHandler.ListClubAssignments)
	}
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}
