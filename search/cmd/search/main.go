package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kust1q/Zapp/search/internal/config"
	searchgrpc "github.com/kust1q/Zapp/search/internal/controllers/grpc/servers/search"
	"github.com/kust1q/Zapp/search/internal/controllers/kafka"
	"github.com/kust1q/Zapp/search/internal/providers/search/elastic"
	"github.com/kust1q/Zapp/search/internal/service/search"
	el "github.com/kust1q/Zapp/search/pkg/elastic"
	searchproto "github.com/kust1q/Zapp/search/pkg/gen/proto/search"
	kafkaProvider "github.com/kust1q/Zapp/search/pkg/kafka"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	if err := config.InitConfig(); err != nil {
		logrus.WithError(err).Fatal("error initializing config")
	}

	cfg := config.Get()
	if err := cfg.Validate(); err != nil {
		logrus.WithError(err).Fatal("invalid configuration")
	}

	elasticClient, err := el.NewElasticClient([]string{fmt.Sprintf("http://%s:%s", cfg.Elastic.Host, cfg.Elastic.Port)})
	if err != nil {
		logrus.Fatalf("failed to connect to elasticsearch: %v", err)
	}
	elasticRepo := elastic.NewElasticRepository(elasticClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	searchService := search.NewSearchService(elasticRepo)

	kafkaHadler := kafka.NewSearchHandler(searchService)
	consumer := kafkaProvider.NewEventConsumer(&cfg.Kafka)

	go func() {
		logrus.Info("Starting Kafka consumer...")
		if taskErr := consumer.Run(ctx, kafkaHadler.Handle); taskErr != nil {
			logrus.Error("kafka consumer failed", taskErr)
		}
	}()

	lis, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%s", cfg.GRPC.SearchPort))
	if err != nil {
		logrus.Fatal("failed to listen", err)
	}

	grpcServer := grpc.NewServer()
	searchHandler := searchgrpc.NewSearchServer(searchService)
	searchproto.RegisterSearchServiceServer(grpcServer, searchHandler)

	reflection.Register(grpcServer)

	go func() {
		logrus.Info("Starting gRPC server on ", cfg.GRPC.SearchPort)
		if err := grpcServer.Serve(lis); err != nil {
			logrus.Fatal("failed to serve", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logrus.Infof("Received signal: %v. Shutting down...", sig)
	cancel()
	grpcServer.GracefulStop()
	time.Sleep(1 * time.Second)

	logrus.Info("Server exited properly")
}
