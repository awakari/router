package model

type Metadata map[string]any

// Uri is a marker type to store the CloudEvents source attribute value type
type Uri string

// UriRef is a marker type to store the CloudEvents source attribute value type
type UriRef string

// KeySubscription represents the id of the matching subscription.
const KeySubscription = "awakarisubscription"
