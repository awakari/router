package service

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
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

func (lm loggingMiddleware) Route(ctx context.Context, msg *event.Event) (err error) {
	defer func() {
		lm.log.Debug(fmt.Sprintf("Route(msg.Id=%s): %s", msg.ID(), err))
	}()
	return lm.svc.Route(ctx, msg)
}
