package matches

import (
	"context"
	"github.com/awakari/router/model"
)

type serviceMock struct {
}

func NewServiceMock() Service {
	return serviceMock{}
}

func (sm serviceMock) Search(ctx context.Context, msgId string, cursor string, limit uint32) (page []model.SubscriptionBase, err error) {
	if msgId == "fail" {
		err = ErrInternal
	} else if cursor == "" {
		page = []model.SubscriptionBase{
			{
				Id: "sub0",
				Destinations: []string{
					"dst0",
				},
			},
			{
				Id: "sub1",
				Destinations: []string{
					"dst1",
					"dst2",
				},
			},
		}
	}
	return
}
