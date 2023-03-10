package service

import (
	"context"
	"github.com/awakari/router/api/grpc/consumer"
	"github.com/awakari/router/api/grpc/matches"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Submit(t *testing.T) {
	svc := NewService(matches.NewServiceMock(), 123, consumer.NewServiceMock())
	cases := map[string]error{
		"ok":      nil,
		"fail":    matches.ErrInternal,
		"missing": consumer.ErrQueueMissing,
		"full":    consumer.ErrQueueFull,
	}
	for k, expectedErr := range cases {
		t.Run(k, func(t *testing.T) {
			msg := event.New()
			msg.SetID(k)
			err := svc.Submit(context.TODO(), &msg)
			assert.ErrorIs(t, err, expectedErr)
		})
	}
}
