package model

type Metadata map[string]any

// Uri is a marker type to store the CloudEvents source attribute value type
type Uri string

// UriRef is a marker type to store the CloudEvents source attribute value type
type UriRef string

// KeyDestination represents the resolved destination route for the message.
const KeyDestination = "awakari_destination"

// KeySubscription represents the id of the matching subscription.
const KeySubscription = "awakari_subscription"
