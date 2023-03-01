package model

type SubscriptionBase struct {

	// Id represents the unique Subscription id.
	Id string

	// Destinations represents a list of target routes associated with the Subscription.
	Destinations []string
}
