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

func (sm serviceMock) SubmitBatch(ctx context.Context, msgs []*event.Event) (count uint32, err error) {
	for _, msg := range msgs {
		switch msg.ID() {
		case "consumer_fail":
			err = consumer.ErrInternal
		case "consumer_queue_missing":
			err = consumer.ErrQueueMissing
		case "matches_fail":
			err = matches.ErrInternal
		case "queue_fail":
			err = queue.ErrInternal
		case "queue_missing":
			err = queue.ErrQueueMissing
		}
		if err != nil {
			break
		}
		count++
	}
	return
}
