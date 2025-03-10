package property

import "context"

type EventStore interface {
	SaveEvent(ctx context.Context, event *Event) error
	GetEventsForFilter(ctx context.Context, filter *EventFilter, limit int, offset int) ([]*Event, error)
	GetMostRecentEventForFilter(ctx context.Context, filter *EventFilter) (*Event, bool, error)
}

type Handler struct {
	store EventStore
}

func NewHandler(store EventStore) *Handler {
	return &Handler{
		store: store,
	}
}

type SortOrder int8

const (
	Ascending SortOrder = iota
	Descending
)

type AmountType int8

const (
	All AmountType = iota
	Expense
	Income
)
