package consumer

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/router/api/grpc/queue"
	format "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	Submit(ctx context.Context, msg *event.Event) (err error)
}

type service struct {
	client ServiceClient
}

var ErrInternal = errors.New("internal failure")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) Submit(ctx context.Context, msg *event.Event) (err error) {
	var msgProto *pb.CloudEvent
	msgProto, err = format.ToProto(msg)
	if err == nil {
		_, err = svc.client.Submit(ctx, msgProto)
		if err != nil {
			err = decodeError(err)
		}
	}
	return
}

func decodeError(src error) (dst error) {
	switch status.Code(src) {
	case codes.OK:
		dst = nil
	case codes.NotFound:
		dst = queue.ErrQueueMissing
	case codes.ResourceExhausted:
		dst = queue.ErrQueueFull
	default:
		dst = fmt.Errorf("%w: %s", ErrInternal, src)
	}
	return
}
