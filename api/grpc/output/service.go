package output

import (
	"context"
	"errors"
	"fmt"
	grpcMsg "github.com/awakari/router/api/grpc/message"
	"github.com/awakari/router/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	Publish(ctx context.Context, msg model.Message) (err error)
}

type service struct {
	client ServiceClient
}

// ErrInternal indicates some unexpected internal failure.
var ErrInternal = errors.New("internal failure")

func NewService(client ServiceClient) Service {
	return service{
		client: client,
	}
}

func (svc service) Publish(ctx context.Context, msg model.Message) (err error) {
	req := encodeMessage(msg)
	_, err = svc.client.Publish(ctx, &req)
	err = decodeError(err)
	return
}

func encodeMessage(src model.Message) (dst grpcMsg.Message) {
	dst.Id = src.Id
	md := make(map[string]*grpcMsg.AttrValue)
	var av *grpcMsg.AttrValue
	for k, v := range src.Metadata {
		switch vt := v.(type) {
		case bool:
			av = &grpcMsg.AttrValue{Attr: &grpcMsg.AttrValue_CeBoolean{CeBoolean: vt}}
		case []byte:
			av = &grpcMsg.AttrValue{Attr: &grpcMsg.AttrValue_CeBytes{CeBytes: vt}}
		case int32:
			av = &grpcMsg.AttrValue{Attr: &grpcMsg.AttrValue_CeInteger{CeInteger: vt}}
		case model.UriRef:
			av = &grpcMsg.AttrValue{Attr: &grpcMsg.AttrValue_CeUriRef{CeUriRef: string(vt)}}
		case model.Uri:
			av = &grpcMsg.AttrValue{Attr: &grpcMsg.AttrValue_CeUri{CeUri: string(vt)}}
		case string:
			av = &grpcMsg.AttrValue{Attr: &grpcMsg.AttrValue_CeString{CeString: vt}}
		case *timestamppb.Timestamp:
			av = &grpcMsg.AttrValue{Attr: &grpcMsg.AttrValue_CeTimestamp{CeTimestamp: vt}}
		default:
			panic(fmt.Sprintf("message encode failure: unrecognzied attribute value type: %T", v))
		}
		md[k] = av
	}
	dst.Metadata = md
	switch dt := src.Data.(type) {
	case *grpcMsg.Message_BinaryData:
		dst.Data = &grpcMsg.Message_BinaryData{BinaryData: dt.BinaryData}
	case *grpcMsg.Message_ProtoData:
		dst.Data = &grpcMsg.Message_ProtoData{ProtoData: dt.ProtoData}
	case *grpcMsg.Message_TextData:
		dst.Data = &grpcMsg.Message_TextData{TextData: dt.TextData}
	default:
		panic(fmt.Sprintf("message encode failure: unrecognzied data type: %T", dt))
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
