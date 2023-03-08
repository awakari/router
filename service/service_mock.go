package service

import (
	"context"
	"github.com/awakari/router/api/grpc/consumer"
	"github.com/awakari/router/api/grpc/matches"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/cloudevents/sdk-go/v2/event"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) Submit(ctx context.Context, msg *event.Event) (err error) {
	switch msg.ID() {
	case "consumer_fail":
		err = consumer.ErrInternal
	case "matches_fail":
		err = matches.ErrInternal
	case "queue_fail":
		err = queue.ErrInternal
	case "queue_full":
		err = queue.ErrQueueFull
	case "queue_missing":
		err = queue.ErrQueueMissing
	}
	return
}
