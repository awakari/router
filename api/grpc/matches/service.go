package matches

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/router/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {

	// Search returns a next page of basic model.Subscription data. The page offset is defined by the cursor.
	// Empty cursor causes the Search to start from the beginning.
	Search(ctx context.Context, msgId string, cursor string, limit uint32) (page []model.SubscriptionBase, err error)
}

type service struct {
	client ServiceClient
}

// ErrInternal indicates some unexpected internal failure.
var ErrInternal = errors.New("matches: internal failure")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) Search(ctx context.Context, msgId string, cursor string, limit uint32) (page []model.SubscriptionBase, err error) {
	req := SearchRequest{
		MsgId:  msgId,
		Limit:  limit,
		Cursor: cursor,
	}
	var resp *SearchResponse
	resp, err = svc.client.Search(ctx, &req)
	if err == nil {
		for _, s := range resp.Page {
			sub := model.SubscriptionBase{
				Id:           s.Id,
				Destinations: s.Dsts,
			}
			page = append(page, sub)
		}
	} else {
		err = decodeError(err)
	}
	return
}

func decodeError(src error) (dst error) {
	switch status.Code(src) {
	case codes.OK:
		dst = nil
	default:
		dst = fmt.Errorf("%w: %s", ErrInternal, src)
	}
	return
}
