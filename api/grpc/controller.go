package grpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/awakari/router/api/grpc/consumer"
	"github.com/awakari/router/api/grpc/matches"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/awakari/router/service"
	format "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/cloudevents/sdk-go/v2/event"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	serviceController struct {
		svc service.Service
	}
)

func NewServiceController(svc service.Service) ServiceServer {
	return serviceController{
		svc: svc,
	}
}

func (sc serviceController) Submit(ctx context.Context, req *pb.CloudEvent) (resp *emptypb.Empty, err error) {
	var msg *event.Event
	msg, err = format.FromProto(req)
	if err == nil {
		err = sc.svc.Submit(ctx, msg)
		err = encodeError(err)
	}
	return &emptypb.Empty{}, err
}

func encodeError(src error) (dst error) {
	switch {
	case src == nil:
		dst = nil
	case errors.Is(src, consumer.ErrInternal):
		dst = status.Error(codes.Internal, fmt.Sprintf("consumer failure: %s", src.Error()))
	case errors.Is(src, matches.ErrInternal):
		dst = status.Error(codes.Internal, fmt.Sprintf("matches failure: %s", src.Error()))
	case errors.Is(src, queue.ErrInternal):
		dst = status.Error(codes.Internal, fmt.Sprintf("queue failure: %s", src.Error()))
	case errors.Is(src, queue.ErrQueueFull):
		dst = status.Error(codes.ResourceExhausted, src.Error())
	case errors.Is(src, queue.ErrQueueMissing):
		dst = status.Error(codes.NotFound, src.Error())
	default:
		dst = status.Error(codes.Internal, src.Error())
	}
	return
}
