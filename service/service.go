package service

import (
	"context"
	"github.com/awakari/router/api/grpc/consumer"
	"github.com/awakari/router/api/grpc/matches"
	"github.com/awakari/router/model"
	"github.com/cloudevents/sdk-go/v2/event"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	Submit(ctx context.Context, msg *event.Event) (err error)
}

type service struct {
	matchesSvc       matches.Service
	matchesBatchSize uint32
	consumerSvc      consumer.Service
}

func NewService(matchesSvc matches.Service, matchesBatchSize uint32, consumerSvc consumer.Service) Service {
	return service{
		matchesSvc:       matchesSvc,
		matchesBatchSize: matchesBatchSize,
		consumerSvc:      consumerSvc,
	}
}

func (svc service) Submit(ctx context.Context, msg *event.Event) (err error) {
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
			sub := m // copy to avoid the data race test error
			g.Go(func() error {
				return svc.routeBySubscription(gCtx, msg, sub)
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
	dsts := sub.Destinations
	if len(dsts) == 1 {
		msgCopy := msg.Clone()
		msgCopy.SetExtension(model.KeySubscription, sub.Id)
		msgCopy.SetExtension(model.KeyDestination, dsts[0])
		// it's expected that many of subscriptions have single route destination only, hence avoid creating a goroutine
		err = svc.consumerSvc.Submit(ctx, &msgCopy)
	} else {
		g, gCtx := errgroup.WithContext(ctx)
		for _, dst := range dsts {
			msgCopy := msg.Clone()
			msgCopy.SetExtension(model.KeySubscription, sub.Id)
			msgCopy.SetExtension(model.KeyDestination, dst)
			g.Go(func() error {
				return svc.consumerSvc.Submit(gCtx, &msgCopy)
			})
		}
		err = g.Wait()
	}
	return
}
