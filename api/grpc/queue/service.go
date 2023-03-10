package queue

import (
	"context"
	"errors"
	"fmt"
	format "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	SetQueue(ctx context.Context, name string, limit uint32) (err error)
	SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error)
	Poll(ctx context.Context, queue string, limit uint32) (msgs []*event.Event, err error)
}

type service struct {
	client ServiceClient
}

var ErrInternal = errors.New("internal failure")

var ErrQueueMissing = errors.New("missing queue")

var ErrQueueFull = errors.New("queue is full")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) SetQueue(ctx context.Context, name string, limit uint32) (err error) {
	req := SetQueueRequest{
		Name:  name,
		Limit: limit,
	}
	_, err = svc.client.SetQueue(ctx, &req)
	if err != nil {
		err = decodeError(err)
	}
	return
}

func (svc service) SubmitMessage(ctx context.Context, queue string, msg *event.Event) (err error) {
	var msgProto *pb.CloudEvent
	msgProto, err = format.ToProto(msg)
	if err == nil {
		req := SubmitMessageRequest{
			Queue: queue,
			Msg:   msgProto,
		}
		_, err = svc.client.SubmitMessage(ctx, &req)
		if err != nil {
			err = decodeError(err)
		}
	}
	return
}

func (svc service) Poll(ctx context.Context, queue string, limit uint32) (msgs []*event.Event, err error) {
	req := PollRequest{
		Queue: queue,
		Limit: limit,
	}
	var resp *PollResponse
	resp, err = svc.client.Poll(ctx, &req)
	if err != nil {
		err = decodeError(err)
	} else {
		var msg *event.Event
		for _, msgProto := range resp.Msgs {
			msg, err = format.FromProto(msgProto)
			if err != nil {
				break
			} else {
				msgs = append(msgs, msg)
			}
		}
	}
	return
}

func decodeError(src error) (dst error) {
	switch status.Code(src) {
	case codes.OK:
		dst = nil
	case codes.NotFound:
		dst = ErrQueueMissing
	case codes.Unavailable:
		dst = ErrQueueFull
	default:
		dst = fmt.Errorf("%w: %s", ErrInternal, src)
	}
	return
}
