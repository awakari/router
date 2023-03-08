package consumer

import (
	"context"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Submit(t *testing.T) {
	svc := NewService(NewClientMock())
	cases := map[string]error{
		"missing": queue.ErrQueueMissing,
		"ok":      nil,
		"full":    queue.ErrQueueFull,
		"fail":    ErrInternal,
	}
	msg := event.New("1.0")
	for k, err := range cases {
		t.Run(k, func(t *testing.T) {
			msg.SetID(k)
			actualErr := svc.Submit(context.TODO(), &msg)
			assert.ErrorIs(t, actualErr, err)
		})
	}
}
