package grpc

import (
	"context"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/awakari/router/service"
	format "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/cloudevents/sdk-go/v2/event"
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

func (sc serviceController) SubmitBatch(ctx context.Context, req *SubmitBatchRequest) (resp *queue.BatchResponse, err error) {
	var msg *event.Event
	var msgs []*event.Event
	for _, msgProto := range req.Msgs {
		msg, err = format.FromProto(msgProto)
		if err != nil {
			break
		}
		msgs = append(msgs, msg)
	}
	resp = &queue.BatchResponse{}
	if err == nil {
		resp.Count, err = sc.svc.SubmitBatch(ctx, msgs)
	}
	if err != nil {
		resp.Err = err.Error()
		err = nil
	}
	return
}
