package output

import (
	"context"
	grpcMsg "github.com/awakari/router/api/grpc/message"
	"github.com/awakari/router/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	clientMock struct {
	}
)

func NewClientMock() ServiceClient {
	return clientMock{}
}

func (cm clientMock) Publish(ctx context.Context, req *grpcMsg.Message, opts ...grpc.CallOption) (resp *emptypb.Empty, err error) {
	if req.Metadata[model.KeyDestination].String() == "fail" {
		err = status.Error(codes.Internal, "")
	}
	resp = &emptypb.Empty{}
	return
}
