package service

import (
	"context"
	"fmt"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/awakari/router/config"
	"github.com/cloudevents/sdk-go/v2/event"
	"strings"
	"time"
)

type queueMiddleware struct {
	svc               Service
	queueSvc          queue.Service
	queueName         string
	queueFallbackName string
	backoffEmpty      time.Duration
	backoffError      time.Duration
	batchSize         uint32
}

func NewQueueMiddleware(svc Service, queueSvc queue.Service, queueConfig config.QueueConfig) Service {
	qm := queueMiddleware{
		svc:               svc,
		queueSvc:          queueSvc,
		queueName:         queueConfig.Name,
		queueFallbackName: fmt.Sprintf("%s-%s", queueConfig.Name, queueConfig.FallBack.Suffix),
		backoffEmpty:      queueConfig.BackoffEmpty,
		backoffError:      queueConfig.BackoffError,
		batchSize:         queueConfig.BatchSize,
	}
	go qm.processQueueLoop()
	return qm
}

func (qm queueMiddleware) SubmitBatch(ctx context.Context, msgs []*event.Event) (count uint32, err error) {
	return qm.queueSvc.SubmitMessageBatch(ctx, qm.queueName, msgs)
}

func (qm queueMiddleware) processQueueLoop() {
	ctx := context.TODO()
	for {
		err := qm.processQueueOnce(ctx)
		if err != nil {
			time.Sleep(qm.backoffError)
		}
	}
}

func (qm queueMiddleware) processQueueOnce(ctx context.Context) (err error) {
	var msgs []*event.Event
	msgs, err = qm.queueSvc.Poll(ctx, qm.queueName, qm.batchSize)
	if err == nil {
		msgCount := uint32(len(msgs))
		if msgCount == 0 {
			time.Sleep(qm.backoffEmpty)
		} else {
			var accepted, n uint32
			for accepted < msgCount {
				n, err = qm.svc.SubmitBatch(ctx, msgs[accepted:])
				accepted += n
				if err != nil {
					if accepted < msgCount {
						qm.submitToFallback(ctx, msgs[accepted:])
					}
					break
				}
			}
		}
	}
	return
}

func (qm queueMiddleware) submitToFallback(ctx context.Context, msgs []*event.Event) {
	accepted, err := qm.queueSvc.SubmitMessageBatch(ctx, qm.queueFallbackName, msgs)
	if err != nil {
		var msgIds []string
		for _, msg := range msgs {
			msgIds = append(msgIds, msg.ID())
		}
		panic("failed to submit the messages to the fallback queue, dropped message ids:\n\t" + strings.Join(msgIds, "\n\t"))
	}
	if accepted < uint32(len(msgs)) {
		var msgIds []string
		for _, msg := range msgs {
			msgIds = append(msgIds, msg.ID())
		}
		panic("fallback queue has not enough capacity, dropped message ids:\n\t" + strings.Join(msgIds, "\n\t"))
	}
}
