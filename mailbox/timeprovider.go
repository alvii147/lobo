package mailbox

import "time"

// TimeProvider represents a component that provides the current timestamp.
type TimeProvider interface {
	Now() time.Time
}
