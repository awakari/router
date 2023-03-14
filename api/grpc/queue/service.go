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
	"strings"
)

type Service interface {
	SetQueue(ctx context.Context, name string, limit uint32) (err error)
	SubmitMessageBatch(ctx context.Context, queue string, msgs []*event.Event) (count uint32, err error)
	Poll(ctx context.Context, queue string, limit uint32) (msgs []*event.Event, err error)
}

type service struct {
	client ServiceClient
}

var ErrInternal = errors.New("queue: internal failure")

var ErrQueueMissing = errors.New("missing queue")

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

func (svc service) SubmitMessageBatch(ctx context.Context, queue string, msgs []*event.Event) (count uint32, err error) {
	var msgProto *pb.CloudEvent
	var msgProtos []*pb.CloudEvent
	for _, msg := range msgs {
		msgProto, err = format.ToProto(msg)
		if err != nil {
			break
		}
		msgProtos = append(msgProtos, msgProto)
	}
	if err == nil {
		req := SubmitMessageBatchRequest{
			Queue: queue,
			Msgs:  msgProtos,
		}
		var resp *BatchResponse
		resp, err = svc.client.SubmitMessageBatch(ctx, &req)
		if err != nil {
			err = decodeError(err)
		} else {
			count = resp.Count
			err = decodeRespError(resp.Err)
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
	default:
		dst = fmt.Errorf("%w: %s", ErrInternal, src)
	}
	return
}

func decodeRespError(src string) (err error) {
	switch {
	case strings.HasPrefix(src, ErrInternal.Error()):
		err = fmt.Errorf("%w: %s", ErrInternal, src[len(ErrInternal.Error()):])
	case strings.HasPrefix(src, ErrQueueMissing.Error()):
		err = fmt.Errorf("%w: %s", ErrQueueMissing, src[len(ErrQueueMissing.Error()):])
	case src == "":
		err = nil
	default:
		err = errors.New(src)
	}
	return
}
