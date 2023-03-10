package service

import (
	"context"
	"fmt"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/awakari/router/config"
	"github.com/cloudevents/sdk-go/v2/event"
	"time"
)

type queueMiddleware struct {
	svc               Service
	queueSvc          queue.Service
	queueName         string
	queueFallbackName string
	sleepOnEmpty      time.Duration
	sleepOnError      time.Duration
	batchSize         uint32
}

func NewQueueMiddleware(svc Service, queueSvc queue.Service, queueConfig config.QueueConfig) Service {
	qm := queueMiddleware{
		svc:               svc,
		queueSvc:          queueSvc,
		queueName:         queueConfig.Name,
		queueFallbackName: fmt.Sprintf("%s-%s", queueConfig.Name, queueConfig.FallBack.Suffix),
		sleepOnEmpty:      time.Duration(queueConfig.SleepOnEmptyMillis) * time.Millisecond,
		sleepOnError:      time.Duration(queueConfig.SleepOnErrorMillis) * time.Millisecond,
		batchSize:         queueConfig.BatchSize,
	}
	go qm.processQueueLoop()
	return qm
}

func (qm queueMiddleware) Submit(ctx context.Context, msg *event.Event) (err error) {
	return qm.queueSvc.SubmitMessage(ctx, qm.queueName, msg)
}

func (qm queueMiddleware) processQueueLoop() {
	ctx := context.TODO()
	for {
		err := qm.processQueueOnce(ctx)
		if err != nil {
			time.Sleep(qm.sleepOnError)
		}
	}
}

func (qm queueMiddleware) processQueueOnce(ctx context.Context) (err error) {
	var msgs []*event.Event
	msgs, err = qm.queueSvc.Poll(ctx, qm.queueName, qm.batchSize)
	if err == nil {
		if len(msgs) == 0 {
			time.Sleep(qm.sleepOnEmpty)
		} else {
			for _, msg := range msgs {
				qm.processMessage(ctx, msg)
			}
		}
	}
	return
}

func (qm queueMiddleware) processMessage(ctx context.Context, msg *event.Event) {
	err := qm.svc.Submit(ctx, msg)
	if err != nil {
		err = qm.queueSvc.SubmitMessage(ctx, qm.queueFallbackName, msg)
	}
	if err != nil {
		panic(err)
	}
}
