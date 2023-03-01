package matches

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	clientMock struct {
	}
)

func NewClientMock() ServiceClient {
	return clientMock{}
}

func (cm clientMock) Search(ctx context.Context, in *SearchRequest, opts ...grpc.CallOption) (resp *SearchResponse, err error) {
	if in.MsgId == "fail" {
		err = status.Error(codes.Internal, "")
	} else {
		resp = &SearchResponse{
			Page: []*SubscriptionOutput{
				{
					Id: "sub0",
					Dsts: []string{
						"dst0",
					},
				},
				{
					Id: "sub1",
					Dsts: []string{
						"dst1",
					},
				},
			},
		}
	}
	return
}
