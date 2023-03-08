package service

import (
	"context"
	"github.com/awakari/router/api/grpc/matches"
	"github.com/awakari/router/api/grpc/queue"
	"github.com/awakari/router/model"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/sync/errgroup"
	"time"
)

type Service interface {
	Route(ctx context.Context, msg *event.Event) (err error)
}

type service struct {
	matchesSvc       matches.Service
	matchesBatchSize uint32
	queueSvc         queue.Service
	queueName        string
	msgBatchSize     uint32
	sleepOnError     time.Duration
	sleepOnEmpty     time.Duration
}

func NewService(matchesSvc matches.Service, matchesBatchSize uint32, queueSvc queue.Service, queueName string, msgBatchSize uint32) Service {
	return service{
		matchesSvc:       matchesSvc,
		matchesBatchSize: matchesBatchSize,
		queueSvc:         queueSvc,
		queueName:        queueName,
		msgBatchSize:     msgBatchSize,
	}
}

func (svc service) Route(ctx context.Context, msg *event.Event) (err error) {
	err = svc.queueSvc.SubmitMessage(ctx, svc.queueName, msg)
	return
}

func (svc service) ProcessLoop() {
	ctx := context.TODO()
	for {
		err := svc.process(ctx)
		if err != nil {
			time.Sleep(svc.sleepOnError)
		}
	}
}

func (svc service) process(ctx context.Context) (err error) {
	var msgs []*event.Event
	msgs, err = svc.queueSvc.Poll(ctx, svc.queueName, svc.msgBatchSize)
	if err == nil {
		if len(msgs) == 0 {
			time.Sleep(svc.sleepOnEmpty)
		} else {
			for _, msg := range msgs {
				err = svc.processMessage(ctx, msg)
				if err != nil {
					break
				}
			}
		}
	}
	return
}

func (svc service) processMessage(ctx context.Context, msg *event.Event) (err error) {
	msgId := msg.ID()
	var cursor string
	var matchesPage []model.SubscriptionBase
	for {
		matchesPage, err = svc.matchesSvc.Search(ctx, msgId, cursor, svc.matchesBatchSize)
		if err != nil || len(matchesPage) == 0 {
			break
		}
		g, gCtx := errgroup.WithContext(ctx)
		for i, m := range matchesPage {
			if i == len(matchesPage)-1 {
				cursor = m.Id
			}
			g.Go(func() error {
				return svc.routeBySubscription(gCtx, msg, m)
			})
		}
		err = g.Wait()
		if err != nil {
			break
		}
	}
	return
}

func (svc service) routeBySubscription(ctx context.Context, msg *event.Event, sub model.SubscriptionBase) (err error) {
	msgCopy := msg.Clone()
	msgCopy.SetExtension(model.KeySubscription, sub.Id)
	dsts := sub.Destinations
	if len(dsts) == 1 {
		// it's expected that many of subscriptions have single route destination only, hence avoid creating a goroutine
		err = svc.queueSvc.SubmitMessage(ctx, dsts[0], &msgCopy)
	} else {
		g, gCtx := errgroup.WithContext(ctx)
		for _, dst := range dsts {
			g.Go(func() error {
				return svc.queueSvc.SubmitMessage(gCtx, dst, &msgCopy)
			})
		}
		err = g.Wait()
	}
	return
}
