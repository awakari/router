package service

import (
	"context"
	"github.com/awakari/router/api/grpc/matches"
	"github.com/awakari/router/api/grpc/output"
	"github.com/awakari/router/model"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	Route(ctx context.Context, msg model.Message) (err error)
}

type service struct {
	matchesSvc       matches.Service
	matchesBatchSize uint32
	outputSvc        output.Service
}

func NewService(matchesSvc matches.Service, matchesBatchSize uint32, outputSvc output.Service) Service {
	return service{
		matchesSvc:       matchesSvc,
		matchesBatchSize: matchesBatchSize,
		outputSvc:        outputSvc,
	}
}

func (svc service) Route(ctx context.Context, msg model.Message) (err error) {
	msgId := msg.Id
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

func (svc service) routeBySubscription(ctx context.Context, msg model.Message, sub model.SubscriptionBase) (err error) {
	md := msg.Metadata
	msg.Metadata = make(map[string]any)
	for k, v := range md {
		msg.Metadata[k] = v
	}
	msg.Metadata[model.KeySubscription] = sub.Id
	dsts := sub.Destinations
	if len(dsts) == 1 {
		// it's expected that many of subscriptions have single route destination only, hence avoid creating a goroutine
		err = svc.routeToDestination(ctx, msg, dsts[0])
	} else {
		g, gCtx := errgroup.WithContext(ctx)
		for _, dst := range dsts {
			g.Go(func() error {
				return svc.routeToDestination(gCtx, msg, dst)
			})
		}
		err = g.Wait()
	}
	return
}

func (svc service) routeToDestination(ctx context.Context, msg model.Message, dst string) (err error) {
	md := msg.Metadata
	msg.Metadata = make(map[string]any)
	for k, v := range md {
		msg.Metadata[k] = v
	}
	msg.Metadata[model.KeyDestination] = dst
	return svc.outputSvc.Publish(ctx, msg)
}
