package main

import (
	"context"
	"fmt"
	grpcApi "github.com/awakari/router/api/grpc"
	"github.com/awakari/router/api/grpc/consumer"
	"github.com/awakari/router/api/grpc/matches"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/awakari/router/config"
	"github.com/awakari/router/service"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func main() {
	//
	slog.Info("starting...")
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		slog.Error("failed to load the config", err)
	}
	opts := slog.HandlerOptions{
		Level: slog.Level(cfg.Log.Level),
	}
	log := slog.New(opts.NewTextHandler(os.Stdout))
	//
	matchesConn, err := grpc.Dial(cfg.Api.Matches.Uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect the matches service", err)
	}
	matchesClient := matches.NewServiceClient(matchesConn)
	matchesSvc := matches.NewService(matchesClient)
	matchesSvc = matches.NewLoggingMiddleware(matchesSvc, log)
	//
	consumerConn, err := grpc.Dial(cfg.Api.Consumer.Uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect the consumer service", err)
	}
	consumerClient := consumer.NewServiceClient(consumerConn)
	consumerSvc := consumer.NewService(consumerClient)
	consumerSvc = consumer.NewLoggingMiddleware(consumerSvc, log)
	//
	queueConn, err := grpc.Dial(cfg.Queue.Uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect the queue service", err)
	}
	queueClient := queue.NewServiceClient(queueConn)
	queueSvc := queue.NewService(queueClient)
	queueSvc = queue.NewLoggingMiddleware(queueSvc, log)
	err = queueSvc.SetQueue(context.TODO(), cfg.Queue.Name, cfg.Queue.Limit)
	if err != nil {
		log.Error("failed to create the work queue", err)
	}
	if cfg.Queue.FallBack.Enabled {
		err = queueSvc.SetQueue(context.TODO(), fmt.Sprintf("%s-%s", cfg.Queue.Name, cfg.Queue.FallBack.Suffix), cfg.Queue.Limit)
	}
	if err != nil {
		log.Error("failed to create the fallback queue", err)
	}
	//
	svc := service.NewService(matchesSvc, cfg.Api.Matches.BatchSize, consumerSvc, cfg.Api.Consumer.Backoff)
	svc = service.NewLoggingMiddleware(svc, log)
	svc = service.NewQueueMiddleware(svc, queueSvc, cfg.Queue)
	//
	log.Info("connected, starting to listen for incoming requests...")
	if err = grpcApi.Serve(svc, cfg.Api.Port); err != nil {
		log.Error("", err)
	}
}
