package events

type InconsistentEvent struct {
	ID        int64
	Direction string // SRC || DST
	Type      string
}

const (
	// InconsistentEventTypeTargetMissing target_missings
	InconsistentEventTypeTargetMissing = "target_missing"
	// InconsistentEventTypeNEQ neq
	InconsistentEventTypeNEQ = "neq"
	// InconsistentEventTypeBaseMissing base_missing
	InconsistentEventTypeBaseMissing = "base_missing"
)
