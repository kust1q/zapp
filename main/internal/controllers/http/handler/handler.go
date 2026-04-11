package http

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	RefreshTokenCookieName = "refresh_token"
)

var (
	// Счетчик HTTP-запросов
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Количество HTTP-запросов",
		},
		[]string{"method", "endpoint", "status"},
	)

	// Гистограмма времени обработки
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Время обработки запроса (секунды)",
			Buckets: []float64{0.1, 0.3, 0.5, 1, 3, 5},
		},
		[]string{"method", "endpoint"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

type Handler struct {
	authService         authService
	tweetService        tweetService
	userService         userService
	clientSearchService clientSearchService
	feedService         feedService
	mediaService        mediaService
	webSocketService    webSocketService
	notificationService notificationService
}

func NewHandler(
	authService authService,
	tweetService tweetService,
	userService userService,
	clientSearchService clientSearchService,
	feedService feedService,
	mediaService mediaService,
	webSocketService webSocketService,
	notificationService notificationService,
) *Handler {
	return &Handler{
		authService:         authService,
		tweetService:        tweetService,
		userService:         userService,
		clientSearchService: clientSearchService,
		feedService:         feedService,
		mediaService:        mediaService,
		webSocketService:    webSocketService,
		notificationService: notificationService,
	}
}

// InitRouters initializes all API routes.
//
// @Router /public [get]
func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	router.Use(h.metricsMiddleware())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowWebSockets = true
	config.AddAllowHeaders("Authorization", "Content-Type")
	config.AllowCredentials = true
	router.Use(cors.New(config))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")

	api.GET("/default", h.getDefaultFeed)

	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.PATCH("/refresh", h.refresh)
		auth.DELETE("/sign-out", h.signOut)
		auth.POST("/forgot-password", h.forgotPassword)
		auth.PATCH("/recovery-password", h.recoveryPassword)
	}

	public := api.Group("/public")
	{
		public.GET("/", h.getDefaultFeed)
		public.GET("/tweets/:tweet_id", h.getTweetById)
		public.GET("/tweets/:tweet_id/replies", h.getReplies)
		public.GET("/tweets/:tweet_id/likes", h.getLikes)
		public.GET("/tweets/media/:tweet_id", h.getTweetMedia)
		public.GET("/users/:username/profile", h.getUserProfile)
		public.GET("/users/:username/tweets", h.getTweetsAndRetweetsByUsername)
		public.GET("/users/:username/followers", h.followers)
		public.GET("/users/:username/following", h.following)
		public.GET("/users/avatar/:user_id", h.getAvatar)
		public.GET("/search", h.search)
	}

	protected := api.Group("/protected", h.authMiddleware)
	{
		protected.GET("/ws", h.serveWs)

		protected.PUT("/reset-password", h.updatePassword)

		tweets := protected.Group("/tweets")
		{
			tweets.POST("", h.createTweet)
			tweets.PATCH("/:tweet_id", h.updateTweet)
			tweets.DELETE("/:tweet_id", h.deleteTweet)

			tweets.POST("/:tweet_id/like", h.likeTweet)
			tweets.DELETE("/:tweet_id/like", h.unlikeTweet)

			tweets.POST("/:tweet_id/reply", h.replyToTweet)

			tweets.POST("/:tweet_id/retweet", h.retweet)
			tweets.DELETE("/:tweet_id/retweet", h.deleteRetweet)

			tweets.DELETE("/:tweet_id/media", h.deleteTweetMedia)
		}

		users := protected.Group("/users")
		{
			users.GET("/me", h.getMe)
			users.PATCH("/me", h.updateMe)
			users.DELETE("/me", h.deleteMe)
			users.POST("/:user_id/follow", h.followUser)
			users.DELETE("/:user_id/follow", h.unfollowUser)
		}
		protected.GET("/feed", h.getFeed)
	}

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint not found"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return router
}
