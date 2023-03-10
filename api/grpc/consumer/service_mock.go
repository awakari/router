package consumer

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) Submit(ctx context.Context, msg *event.Event) (err error) {
	switch msg.ID() {
	case "missing":
		err = ErrQueueMissing
	case "full":
		err = ErrQueueFull
	case "fail":
		err = ErrInternal
	}
	return
}
