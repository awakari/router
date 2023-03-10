package service

import (
	"context"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/awakari/router/config"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueueMiddleware_Submit(t *testing.T) {
	qm := NewQueueMiddleware(NewServiceMock(), queue.NewServiceMock([]*event.Event{}), config.QueueConfig{})
	cases := map[string]error{
		"ok":      nil,
		"fail":    queue.ErrInternal,
		"full":    queue.ErrQueueFull,
		"missing": queue.ErrQueueMissing,
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			msg := event.New()
			msg.SetID(k)
			err := qm.Submit(context.TODO(), &msg)
			assert.ErrorIs(t, err, c)
		})
	}
}
