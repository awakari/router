package main

import (
	grpcApi "github.com/awakari/router/api/grpc"
	"github.com/awakari/router/api/grpc/matches"
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
		Level: cfg.Log.Level,
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
	outputConn, err := grpc.Dial(cfg.Api.Output.Uri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect the output service", err)
	}
	outputClient := output.NewServiceClient(outputConn)
	outputSvc := output.NewService(outputClient)
	outputSvc = output.NewLoggingMiddleware(outputSvc, log)
	//
	svc := service.NewService(matchesSvc, cfg.Api.Matches.BatchSize, outputSvc)
	svc = service.NewLoggingMiddleware(svc, log)
	log.Info("connected, starting to listen for incoming requests...")
	if err = grpcApi.Serve(svc, cfg.Api.Port); err != nil {
		log.Error("", err)
	}
}
