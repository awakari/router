package grpc

import (
	"context"
	"fmt"
	grpcMsg "github.com/awakari/router/api/grpc/message"
	"github.com/awakari/router/model"
	"github.com/awakari/router/service"
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

func (sc serviceController) Route(ctx context.Context, req *grpcMsg.Message) (*emptypb.Empty, error) {
	msg := decodeMessage(req)
	err := sc.svc.Route(ctx, msg)
	err = encodeError(err)
	return &emptypb.Empty{}, err
}

func decodeMessage(src *grpcMsg.Message) (msg model.Message) {
	msg.Id = src.Id
	md := make(map[string]any)
	for k, v := range src.Metadata {
		switch at := v.Attr.(type) {
		case *grpcMsg.AttrValue_CeBoolean:
			md[k] = at.CeBoolean
		case *grpcMsg.AttrValue_CeBytes:
			md[k] = at.CeBytes
		case *grpcMsg.AttrValue_CeInteger:
			md[k] = at.CeInteger
		case *grpcMsg.AttrValue_CeString:
			md[k] = at.CeString
		case *grpcMsg.AttrValue_CeTimestamp:
			md[k] = at.CeTimestamp
		case *grpcMsg.AttrValue_CeUri:
			md[k] = model.Uri(at.CeUri)
		case *grpcMsg.AttrValue_CeUriRef:
			md[k] = model.UriRef(at.CeUriRef)
		default:
			panic(fmt.Sprintf("message decode failure: unrecognzied attribute value type: %T", v))
		}
	}
	msg.Metadata = md
	msg.Data = src.Data
	return
}

func encodeError(svcErr error) (err error) {
	switch {
	case svcErr == nil:
		err = nil
	default:
		err = status.Error(codes.Internal, svcErr.Error())
	}
	return
}
