package matches

import (
	"context"
	"github.com/awakari/router/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Search(t *testing.T) {
	svc := NewService(NewClientMock())
	cases := map[string]struct {
		page []model.SubscriptionBase
		err  error
	}{
		"ok": {
			page: []model.SubscriptionBase{
				{
					Id:           "sub0",
					Destinations: []string{"dst0"},
				},
				{
					Id:           "sub1",
					Destinations: []string{"dst1"},
				},
			},
		},
		"fail": {
			err: ErrInternal,
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			page, err := svc.Search(context.TODO(), k, "", 0)
			assert.ErrorIs(t, err, c.err)
			if c.err == nil {
				assert.Equal(t, c.page, page)
			}
		})
	}
}
