package queue

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

var MsgProtoMocks = []*pb.CloudEvent{
	{
		Id:          "3426d090-1b8a-4a09-ac9c-41f2de24d5ac",
		Source:      "source0",
		SpecVersion: "1.0",
		Type:        "type0",
		Attributes: map[string]*pb.CloudEventAttributeValue{
			"foo": {
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: "bar",
				},
			},
			"datacontenttype": {
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: "text/plain",
				},
			},
		},
		Data: &pb.CloudEvent_BinaryData{
			BinaryData: []byte("yohoho"),
		},
	},
	{
		Id:          "f7102c87-3ce4-4bb0-8527-b4644f685b13",
		Source:      "source1",
		SpecVersion: "1.0",
		Type:        "type1",
		Attributes: map[string]*pb.CloudEventAttributeValue{
			"bool": {
				Attr: &pb.CloudEventAttributeValue_CeBoolean{
					CeBoolean: true,
				},
			},
			"datacontenttype": {
				Attr: &pb.CloudEventAttributeValue_CeString{
					CeString: "application/octet-stream",
				},
			},
		},
		Data: &pb.CloudEvent_BinaryData{
			BinaryData: []byte{1, 2, 3},
		},
	},
}

func (cm clientMock) SetQueue(ctx context.Context, in *SetQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	var err error
	if in.Name == "fail" {
		err = status.Error(codes.Internal, "failed to set up the queue")
	}
	return &emptypb.Empty{}, err
}

func (cm clientMock) SubmitMessage(ctx context.Context, in *SubmitMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	var err error
	switch in.Queue {
	case "fail":
		err = status.Error(codes.Internal, "failed to submit the message")
	case "full":
		err = status.Error(codes.ResourceExhausted, "failed to submit the message")
	case "missing":
		err = status.Error(codes.NotFound, "failed to submit the message")
	}
	return &emptypb.Empty{}, err
}

func (cm clientMock) Poll(ctx context.Context, in *PollRequest, opts ...grpc.CallOption) (*PollResponse, error) {
	resp := &PollResponse{}
	var err error
	switch in.Queue {
	case "fail":
		err = status.Error(codes.Internal, "failed to poll")
	case "missing":
		err = status.Error(codes.NotFound, "failed to poll")
	default:
		resp.Msgs = MsgProtoMocks
	}
	return resp, err
}
