package grpc

import (
	"context"
	"fmt"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/awakari/router/service"
	"github.com/cloudevents/sdk-go/binding/format/protobuf/v2/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	cases := map[string]struct {
		msgIds []string
		resp   *queue.BatchResponse
		err    error
	}{
		"ok": {
			msgIds: []string{
				"msg0",
				"msg1",
				"msg2",
			},
			resp: &queue.BatchResponse{
				Count: 3,
				Err:   "",
			},
		},
		"consumer_fail": {
			msgIds: []string{
				"consumer_fail",
				"msg1",
				"msg2",
			},
			resp: &queue.BatchResponse{
				Count: 0,
				Err:   "consumer: internal failure",
			},
		},
		"consumer_queue_missing": {
			msgIds: []string{
				"msg0",
				"consumer_queue_missing",
				"msg2",
			},
			resp: &queue.BatchResponse{
				Count: 1,
				Err:   "consumer: missing queue",
			},
		},
		"matches_fail": {
			msgIds: []string{
				"msg0",
				"msg1",
				"matches_fail",
			},
			resp: &queue.BatchResponse{
				Count: 2,
				Err:   "matches: internal failure",
			},
		},
		"queue_fail": {
			msgIds: []string{
				"queue_fail",
				"msg1",
				"msg2",
			},
			resp: &queue.BatchResponse{
				Count: 0,
				Err:   "queue: internal failure",
			},
		},
		"queue_missing": {
			msgIds: []string{
				"msg0",
				"queue_missing",
				"msg2",
			},
			resp: &queue.BatchResponse{
				Count: 1,
				Err:   "missing queue",
			},
		},
	}
	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			var msgs []*pb.CloudEvent
			for _, msgId := range c.msgIds {
				msgs = append(msgs, &pb.CloudEvent{Id: msgId})
			}
			resp, err := client.SubmitBatch(context.TODO(), &SubmitBatchRequest{Msgs: msgs})
			assert.Equal(t, c.resp.Count, resp.Count)
			assert.Equal(t, c.resp.Err, resp.Err)
			assert.Nil(t, err)
		})
	}
}
