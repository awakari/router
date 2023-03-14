package service

import (
	"context"
	"github.com/awakari/router/api/grpc/consumer"
	"github.com/awakari/router/api/grpc/matches"
	"github.com/awakari/router/model"
	"github.com/cloudevents/sdk-go/v2/event"
	"time"
)

type Service interface {
	SubmitBatch(ctx context.Context, msgs []*event.Event) (count uint32, err error)
}

type service struct {
	matchesSvc       matches.Service
	matchesBatchSize uint32
	consumerSvc      consumer.Service
	consumerBackoff  time.Duration
}

func NewService(matchesSvc matches.Service, matchesBatchSize uint32, consumerSvc consumer.Service, consumerBackoff time.Duration) Service {
	return service{
		matchesSvc:       matchesSvc,
		matchesBatchSize: matchesBatchSize,
		consumerSvc:      consumerSvc,
		consumerBackoff:  consumerBackoff,
	}
}

func (svc service) SubmitBatch(ctx context.Context, msgs []*event.Event) (count uint32, err error) {
	for _, msg := range msgs {
		err = svc.processMessage(ctx, msg)
		if err != nil {
			break
		}
		count++
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
		cursor = matchesPage[len(matchesPage)-1].Id
		dstMsgs := assignDestinations(matchesPage, msg)
		dstMsgCount := uint32(len(dstMsgs))
		var consumed, n uint32
		for consumed < dstMsgCount {
			n, err = svc.consumerSvc.SubmitBatch(ctx, dstMsgs[consumed:])
			consumed += n
			if err != nil {
				break
			}
			if n == 0 {
				time.Sleep(svc.consumerBackoff)
			}
		}
		if err != nil {
			break
		}
	}
	return
}

func assignDestinations(matches []model.SubscriptionBase, srcMsg *event.Event) (dstMsgs []*event.Event) {
	for _, sub := range matches {
		msgSubAssigned := srcMsg.Clone()
		msgSubAssigned.SetExtension(model.KeySubscription, sub.Id)
		for _, dst := range sub.Destinations {
			msgDstAssigned := msgSubAssigned.Clone()
			msgDstAssigned.SetExtension(model.KeyDestination, dst)
			dstMsgs = append(dstMsgs, &msgDstAssigned)
		}
	}
	return
}
