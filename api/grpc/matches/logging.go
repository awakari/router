package matches

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

func (lm loggingMiddleware) Search(ctx context.Context, msgId string, cursor string, limit uint32) (page []model.SubscriptionBase, err error) {
	defer func() {
		lm.log.Debug(fmt.Sprintf("matches.Search(msgId=%s, cursor=%s, limit=%d): %d, %s", msgId, cursor, limit, len(page), err))
	}()
	return lm.svc.Search(ctx, msgId, cursor, limit)
}
