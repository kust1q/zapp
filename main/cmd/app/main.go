// @title           Twitter-like API
// @version         1.0
// @description     Simple Twitter-like REST API.
// @BasePath        /api/v1/

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description  Use "Bearer {access_token}" format.
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/kust1q/Zapp/main/docs"
	"github.com/kust1q/Zapp/main/internal/config"
	tweetgrpc "github.com/kust1q/Zapp/main/internal/controllers/grpc/servers/tweet"
	usergrpc "github.com/kust1q/Zapp/main/internal/controllers/grpc/servers/user"
	httpHandler "github.com/kust1q/Zapp/main/internal/controllers/http/handler"
	s3 "github.com/kust1q/Zapp/main/internal/providers/db/minio"
	db "github.com/kust1q/Zapp/main/internal/providers/db/postgres"
	"github.com/kust1q/Zapp/main/internal/providers/db/redis/cache"
	"github.com/kust1q/Zapp/main/internal/providers/db/redis/tokens"
	searchClient "github.com/kust1q/Zapp/main/internal/providers/search"
	wsProvider "github.com/kust1q/Zapp/main/internal/providers/websocket" // Infrastructure
	"github.com/kust1q/Zapp/main/internal/service/auth"
	"github.com/kust1q/Zapp/main/internal/service/feed"
	"github.com/kust1q/Zapp/main/internal/service/media"
	"github.com/kust1q/Zapp/main/internal/service/notification"
	searchService "github.com/kust1q/Zapp/main/internal/service/search"
	"github.com/kust1q/Zapp/main/internal/service/tweets"
	"github.com/kust1q/Zapp/main/internal/service/user" // Connection Logic
	"github.com/kust1q/Zapp/main/internal/service/websocket"
	"github.com/kust1q/Zapp/main/internal/domain/entity"
	tweetproto "github.com/kust1q/Zapp/main/pkg/gen/proto/tweet"
	userproto "github.com/kust1q/Zapp/main/pkg/gen/proto/user"
	kafkaProvider "github.com/kust1q/Zapp/main/pkg/kafka"
	mn "github.com/kust1q/Zapp/main/pkg/minio"
	pg "github.com/kust1q/Zapp/main/pkg/postgres"
	rs "github.com/kust1q/Zapp/main/pkg/redis"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// --- Configs ---
	if err := config.InitConfig(); err != nil {
		logrus.WithError(err).Fatal("error initializing config")
	}
	cfg := config.Get()
	if err := cfg.Validate(); err != nil {
		logrus.WithError(err).Fatal("invalid configuration")
	}

	redisClient, err := rs.NewRedisClient(&cfg.Redis)
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize redis client")
	}
	cache := cache.NewCache(redisClient, cfg.Cache.DefaultTtl, cfg.Cache.CountersTtl)

	postgresConnect, err := pg.NewPostgresConnect(&cfg.Postgres)
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize pg connect")
	}
	pgDB := db.NewPostgresDB(postgresConnect, cache)

	minioClient, err := mn.NewMinioClient(&cfg.Minio)
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize minio client")
	}

	mediaTypeMap := map[entity.MediaType]entity.MediaPolicy{
		entity.MediaTypeImage: {
			MaxSize:     16 * 1024 * 1024,
			AllowedMime: []string{"image/jpeg", "image/png", "image/webp"},
			AllowedExt:  []string{".jpg", ".jpeg", ".png", ".webp", ".bmp", ".tiff"},
		},
		entity.MediaTypeVideo: {
			MaxSize:     512 * 1024 * 1024,
			AllowedMime: []string{"video/mp4", "video/quicktime", "video/x-m4v"},
			AllowedExt:  []string{".mp4", ".mov", ".m4v", ".avi", ".wmv", ".flv", ".webm"},
		},
		entity.MediaTypeGIF: {
			MaxSize:       16 * 1024 * 1024,
			AllowedMime:   []string{"image/gif"},
			AllowedExt:    []string{".gif"},
			ForceMimeType: "image/gif",
		},
		entity.MediaTypeAudio: {
			MaxSize:     64 * 1024 * 1024,
			AllowedMime: []string{"audio/mpeg", "audio/wav", "audio/x-wav", "audio/ogg", "audio/flac", "audio/aac", "audio/x-m4a", "audio/webm"},
			AllowedExt:  []string{".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a"},
		},
	}
	minioDB := s3.NewMinioDB(minioClient, &cfg.Minio, mediaTypeMap)

	kafkaProducer := kafkaProvider.NewEventProducer(&cfg.Kafka)
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			logrus.Errorf("Error closing Kafka producer: %v", err)
		}
	}()

	searchConn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", cfg.GRPC.Host, cfg.GRPC.SearchPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrus.Fatalf("failed to create grpc client: %v", err)
	}
	defer searchConn.Close()

	searchClient := searchClient.NewClientSearchService(searchConn)
	defer searchConn.Close()

	wsHub := wsProvider.NewHub()
	go wsHub.Run()

	mediaService := media.NewMediaService(pgDB, minioDB)
	authService := auth.NewAuthService(
		&config.AuthServiceConfig{PrivateKey: cfg.JWT.PrivateKey, PublicKey: cfg.JWT.PublicKey, AccessTTL: cfg.Tokens.AccessTTL, RefreshTTL: cfg.Tokens.RefreshTTL},
		pgDB,
		mediaService,
		tokens.NewTokenStorage(redisClient),
		kafkaProducer)
	tweetService := tweets.NewTweetService(pgDB, mediaService, kafkaProducer)
	userService := user.NewUserService(pgDB, mediaService, kafkaProducer)
	feedService := feed.NewFeedService(pgDB, tweetService)
	searchService := searchService.NewSearchService(pgDB, mediaService, tweetService, searchClient)
	wsService := websocket.NewWebSocketService(wsHub)
	notifService := notification.NewNotificationService(wsHub, pgDB)

	handler := httpHandler.NewHandler(
		authService,
		tweetService,
		userService,
		searchService,
		feedService,
		mediaService,
		wsService,
		notifService,
	)

	srv := &http.Server{
		Addr:    ":" + cfg.App.Port,
		Handler: handler.InitRouters(),
	}

	go func() {
		logrus.Infof("Starting server on port %s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	lis, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%s", cfg.GRPC.IntegrationPort))
	if err != nil {
		logrus.Fatalf("failed to listen for grpc: %v", err)
	}
	grpcServer := grpc.NewServer()
	tweetGrpcHandler := tweetgrpc.NewTweetServer(tweetService)
	userGrpcHandler := usergrpc.NewUserServer(userService)

	reflection.Register(grpcServer)

	tweetproto.RegisterTweetServiceServer(grpcServer, tweetGrpcHandler)
	userproto.RegisterUserServiceServer(grpcServer, userGrpcHandler)

	go func() {
		logrus.Infof("Starting Integration gRPC server on port %s", cfg.GRPC.IntegrationPort)
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Fatalf("grpc serve failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server forced to shutdown:", err)
	}

	logrus.Info("Closing database connection...")
	if err := postgresConnect.Close(); err != nil {
		logrus.Errorf("Error closing DB: %v", err)
	}

	logrus.Info("Closing redis connection...")
	if err := redisClient.Close(); err != nil {
		logrus.Errorf("Error closing Redis: %v", err)
	}

	logrus.Info("Server exiting")
}
