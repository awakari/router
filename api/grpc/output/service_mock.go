package output

import (
	"context"
	"github.com/awakari/router/model"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) Publish(ctx context.Context, msg model.Message) (err error) {
	if msg.Id == "fail" {
		err = ErrInternal
	}
	return
}
