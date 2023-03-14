package queue

import (
	"context"
	format "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_SetQueue(t *testing.T) {
	svc := NewService(NewClientMock())
	cases := map[string]error{
		"ok":       nil,
		"existing": nil,
		"fail":     ErrInternal,
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			err := svc.SetQueue(context.TODO(), k, 10)
			assert.ErrorIs(t, err, c)
		})
	}
}

func TestService_SubmitMessageBatch(t *testing.T) {
	svc := NewService(NewClientMock())
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
			err:   ErrInternal,
		},
		"not enough space in the queue": {
			msgIds: []string{
				"msg0",
				"msg1",
				"full",
			},
			count: 2,
		},
		"queue lost": {
			msgIds: []string{
				"missing",
				"msg1",
				"msg2",
			},
			count: 0,
			err:   ErrQueueMissing,
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
			count, err := svc.SubmitMessageBatch(context.TODO(), "queue0", msgs)
			assert.Equal(t, c.count, count)
			assert.ErrorIs(t, err, c.err)
		})
	}
}

func TestService_Poll(t *testing.T) {
	svc := NewService(NewClientMock())
	var msgMocks []*event.Event
	for _, msgProto := range MsgProtoMocks {
		msg, _ := format.FromProto(msgProto)
		msgMocks = append(msgMocks, msg)
	}
	cases := map[string]struct {
		msgs []*event.Event
		err  error
	}{
		"fail": {
			err: ErrInternal,
		},
		"missing": {
			err: ErrQueueMissing,
		},
		"queue0": {
			msgs: msgMocks,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			msgs, err := svc.Poll(context.TODO(), k, 0)
			assert.ErrorIs(t, err, c.err)
			if err == nil {
				assert.Equal(t, c.msgs, msgs)
			}
		})
	}
}
