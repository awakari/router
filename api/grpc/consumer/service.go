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
	"strings"
)

type Service interface {
	SubmitBatch(ctx context.Context, msgs []*event.Event) (count uint32, err error)
}

type service struct {
	client ServiceClient
}

var ErrInternal = errors.New("consumer: internal failure")

var ErrQueueMissing = errors.New("consumer: missing queue")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) SubmitBatch(ctx context.Context, msgs []*event.Event) (count uint32, err error) {
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
		req := SubmitBatchRequest{
			Msgs: msgProtos,
		}
		var resp *queue.BatchResponse
		resp, err = svc.client.SubmitBatch(ctx, &req)
		if err != nil {
			err = decodeError(err)
		} else {
			count = resp.Count
			err = decodeRespError(resp.Err)
		}
	}
	return
}

func decodeError(src error) (dst error) {
	switch status.Code(src) {
	case codes.OK:
		dst = nil
	default:
		dst = fmt.Errorf("%w: consumer: %s", ErrInternal, src)
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
