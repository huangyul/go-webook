package events

import "context"

type Producer interface {
	ProduceEvent(ctx context.Context, event InconsistentEvent) error
}
