package property

import (
	"context"
	"fmt"
	"time"
)

func (h *Handler) GetMonthlyReport(ctx context.Context, PropertyID string, month time.Month, year int, offset int, limit int) ([]*Event, float64, error) {
	if PropertyID == "" {
		return nil, 0, fmt.Errorf("empty property ID")
	}

	// Calculate the start and end of the month
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	startingBalance, err := h.getBalanceForDate(ctx, PropertyID, startOfMonth)
	if err != nil {
		return nil, 0, fmt.Errorf("get balance for date: %v", err)
	}

	events, err := h.GetPropertyEvents(ctx, PropertyID, startOfMonth, endOfMonth, Ascending, All, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("get property events: %v", err)
	}

	return events, startingBalance, nil
}
