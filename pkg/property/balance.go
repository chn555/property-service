package property

import (
	"context"
	"fmt"
	"time"
)

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

func (h *Handler) getBalanceForDate(ctx context.Context, PropertyID string, date time.Time) (float64, error) {
	startingBalanceFilter := &EventFilter{
		PropertyID: PropertyID,
		BeforeTime: date.Add(-1 * time.Millisecond),
	}
	startingBalance, exist, err := h.store.GetMostRecentEventForFilter(ctx, startingBalanceFilter)
	if err != nil {
		return 0, fmt.Errorf("get events for filter: %v", err)
	}
	if !exist {
		return 0, nil
	}
	return startingBalance.PostEventBalance, nil
}
