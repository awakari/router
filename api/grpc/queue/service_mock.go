package queue

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
)

type serviceMock struct {
	msgs []*event.Event
}

func NewServiceMock(msgs []*event.Event) Service {
	return serviceMock{
		msgs: msgs,
	}
}

func (sm serviceMock) SetQueue(ctx context.Context, queue string, limit uint32) (err error) {
	switch queue {
	case "fail":
		err = ErrInternal
	}
	return
}

func (sm serviceMock) SubmitMessageBatch(ctx context.Context, queue string, msgs []*event.Event) (count uint32, err error) {
	for _, msg := range msgs {
		if msg.ID() == "missing" {
			err = ErrQueueMissing
			break
		}
		if msg.ID() == "fail" {
			err = ErrInternal
			break
		}
		if msg.ID() == "full" {
			break
		}
		count++
	}
	return
}

func (sm serviceMock) Poll(ctx context.Context, queue string, limit uint32) (msgs []*event.Event, err error) {
	switch queue {
	case "fail":
		err = ErrInternal
	case "missing":
		err = ErrQueueMissing
	default:
		msgs = sm.msgs
	}
	return
}
