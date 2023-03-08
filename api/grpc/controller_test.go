package grpc

import (
	"context"
	"fmt"
	"github.com/awakari/router/service"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"os"
	"testing"
)

const port = 8080

var log = slog.Default()

func TestMain(m *testing.M) {
	svc := service.NewServiceMock()
	svc = service.NewLoggingMiddleware(svc, log)
	go func() {
		err := Serve(svc, port)
		if err != nil {
			log.Error("", err)
		}
	}()
	code := m.Run()
	os.Exit(code)
}

func TestServiceController_Submit(t *testing.T) {
	//
	addr := fmt.Sprintf("localhost:%d", port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, err)
	client := NewServiceClient(conn)
	//
	cases := map[string]error{
		"ok":            nil,
		"consumer_fail": status.Error(codes.Internal, "consumer failure: internal failure"),
		"matches_fail":  status.Error(codes.Internal, "matches failure: internal failure"),
		"queue_fail":    status.Error(codes.Internal, "queue failure: internal failure"),
		"queue_full":    status.Error(codes.ResourceExhausted, "queue is full"),
		"queue_missing": status.Error(codes.NotFound, "missing queue"),
	}
	//
	for k, expectedErr := range cases {
		t.Run(k, func(t *testing.T) {
			_, err := client.Submit(context.TODO(), &pb.CloudEvent{
				Id: k,
			})
			assert.ErrorIs(t, err, expectedErr)
		})
	}
}
