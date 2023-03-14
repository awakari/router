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
	svc := NewService(matches.NewServiceMock(), 123, consumer.NewServiceMock(), 456)
	cases := map[string]struct {
		msgIds []string
		count  uint32
		err    error
	}{
		"ok": {
			msgIds: []string{
				"msg0",
				"msg1",
				"msg2",
			},
			count: 3,
		},
		"fail on 2nd": {
			msgIds: []string{
				"msg0",
				"fail",
				"msg2",
			},
			count: 1,
			err:   matches.ErrInternal,
		},
		"not enough space in the queue - retries until all consumed": {
			msgIds: []string{
				"msg0",
				"msg1",
				"full",
			},
			count: 3,
		},
		"queue lost": {
			msgIds: []string{
				"missing",
				"msg1",
				"msg2",
			},
			count: 0,
			err:   consumer.ErrQueueMissing,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			var msgs []*event.Event
			for _, msgId := range c.msgIds {
				msg := event.New()
				msg.SetID(msgId)
				msgs = append(msgs, &msg)
			}
			count, err := svc.SubmitBatch(context.TODO(), msgs)
			assert.Equal(t, c.count, count)
			assert.ErrorIs(t, err, c.err)
		})
	}
}
