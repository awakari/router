package consumer

import (
	"context"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type clientMock struct {
}

func NewClientMock() ServiceClient {
	return clientMock{}
}

func (cm clientMock) Submit(ctx context.Context, in *pb.CloudEvent, opts ...grpc.CallOption) (resp *emptypb.Empty, err error) {
	switch in.Id {
	case "missing":
		err = status.Error(codes.NotFound, "destination queue not found")
	case "full":
		err = status.Error(codes.ResourceExhausted, "destination queue is full")
	case "fail":
		err = status.Error(codes.Internal, "internal failure")
	}
	return &emptypb.Empty{}, err
}
