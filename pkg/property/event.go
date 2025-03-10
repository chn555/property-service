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

func (h *Handler) SaveEvent(ctx context.Context, PropertyID string, amount float64, date time.Time) (float64, error) {
	if PropertyID == "" {
		return 0, fmt.Errorf("empty property ID")
	} else if amount == 0 {
		return 0, nil
	} else if date.IsZero() {
		return 0, fmt.Errorf("invalid date")
	}

	curBalance, err := h.GetBalance(ctx, PropertyID)
	if err != nil {
		return 0, fmt.Errorf("get balance: %v", err)
	}

	event := &Event{
		PropertyID:       PropertyID,
		EventAmount:      amount,
		PostEventBalance: curBalance + amount,
		Date:             date,
	}
	if err := h.store.SaveEvent(ctx, event); err != nil {
		return 0, fmt.Errorf("save event: %v", err)
	}
	return event.PostEventBalance, nil
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

	sortByOrder(events, sortOrder)

	return events, nil
}

func sortByOrder(events []*Event, sortOrder SortOrder) {
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
}
