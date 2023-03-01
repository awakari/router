package service

import (
	"context"
	"fmt"
	"github.com/awakari/router/model"
	"golang.org/x/exp/slog"
)

type (
	loggingMiddleware struct {
		svc Service
		log *slog.Logger
	}
)

func NewLoggingMiddleware(svc Service, log *slog.Logger) Service {
	return loggingMiddleware{
		svc: svc,
		log: log,
	}
}

func (lm loggingMiddleware) Route(ctx context.Context, msg model.Message) (err error) {
	defer func() {
		lm.log.Debug(fmt.Sprintf("Route(msg.Id=%s): %s", msg.Id, err))
	}()
	return lm.svc.Route(ctx, msg)
}
