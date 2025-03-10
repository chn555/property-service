package property

import (
	"context"
	"fmt"

	"slices"
	"time"
)

type Event struct {
	PropertyID       string    `json:"property_id,omitempty" bson:"property_id"`
	EventAmount      float64   `json:"event_amount" bson:"event_amount"`
	PostEventBalance float64   `json:"post_event_balance" bson:"post_event_balance"`
	Date             time.Time `json:"date" bson:"date"`
}

type EventFilter struct {
	PropertyID string
	AfterTime  time.Time
	BeforeTime time.Time
	AmountType AmountType
}

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

func (h *Handler) SaveEvent(ctx context.Context, PropertyID string, amount float64, date time.Time) (float64, error) {
	if PropertyID == "" {
		return 0, fmt.Errorf("empty property ID")
	} else if amount == 0 {
		return 0, nil
	} else if date.IsZero() {
		return 0, fmt.Errorf("invalid date")
	}

	filter := &EventFilter{
		PropertyID: PropertyID,
	}

	state, exists, err := h.store.GetMostRecentEventForFilter(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("get events for filter: %v", err)
	}
	if !exists {
		state = &Event{
			PropertyID:       PropertyID,
			EventAmount:      0,
			PostEventBalance: 0,
			Date:             date,
		}
	}

	balance := state.PostEventBalance + amount
	event := &Event{
		PropertyID:       PropertyID,
		EventAmount:      amount,
		PostEventBalance: balance,
		Date:             date,
	}
	if err := h.store.SaveEvent(ctx, event); err != nil {
		return 0, fmt.Errorf("save event: %v", err)
	}
	return event.PostEventBalance, nil
}

func (h *Handler) GetBalance(ctx context.Context, PropertyID string) (float64, error) {
	if PropertyID == "" {
		return 0, fmt.Errorf("empty property ID")
	}

	state, exists, err := h.store.GetMostRecentEventForFilter(ctx, &EventFilter{PropertyID: PropertyID})
	if err != nil {
		return 0, fmt.Errorf("get events for filter: %v", err)
	}
	if !exists {
		return 0, nil
	}
	return state.PostEventBalance, nil
}

func (h *Handler) GetPropertyEvents(ctx context.Context, PropertyID string, dateFrom time.Time, dateTo time.Time, sortOrder SortOrder, amountType AmountType, offset int, limit int) ([]*Event, error) {
	if PropertyID == "" {
		return nil, fmt.Errorf("empty property ID")
	} else if dateFrom.After(dateTo) {
		return nil, fmt.Errorf("dateFrom must be before dateTo")
	}

	filter := &EventFilter{
		PropertyID: PropertyID,
		AfterTime:  dateFrom,
		BeforeTime: dateTo,
		AmountType: amountType,
	}
	events, err := h.store.GetEventsForFilter(ctx, filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get events for filter: %v", err)
	}

	switch sortOrder {
	case Ascending:
		slices.SortFunc(events, func(a, b *Event) int {
			if a.Date.Before(b.Date) {
				return -1
			} else if a.Date.After(b.Date) {
				return 1
			}
			return 0
		})
	case Descending:
		slices.SortFunc(events, func(a, b *Event) int {
			if a.Date.After(b.Date) {
				return -1
			} else if a.Date.Before(b.Date) {
				return 1
			}
			return 0
		})
	}

	return events, nil
}

func (h *Handler) GetMonthlyReport(ctx context.Context, PropertyID string, month time.Month, year int, offset int, limit int) ([]*Event, float64, error) {
	if PropertyID == "" {
		return nil, 0, fmt.Errorf("empty property ID")
	}

	// Calculate the start and end of the month
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	startingBalanceFilter := &EventFilter{
		PropertyID: PropertyID,
		BeforeTime: startOfMonth.Add(-1 * time.Millisecond),
	}
	startingBalance, exist, err := h.store.GetMostRecentEventForFilter(ctx, startingBalanceFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("get events for filter: %v", err)
	}
	if !exist {
		startingBalance = &Event{
			PropertyID:       PropertyID,
			EventAmount:      0,
			PostEventBalance: 0,
			Date:             startOfMonth,
		}
	}

	events, err := h.GetPropertyEvents(ctx, PropertyID, startOfMonth, endOfMonth, Ascending, All, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("get property events: %v", err)
	}

	return events, startingBalance.PostEventBalance, nil
}
