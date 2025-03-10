package property

import (
	"context"
	"github.com/chn555/property-service/internal/rest"
	"github.com/chn555/property-service/pkg/property"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
	"time"
)

type MonthlyReportEvent struct {
	PropertyID  string    `json:"property_id,omitempty" bson:"property_id"`
	EventAmount float64   `json:"event_amount" bson:"event_amount"`
	Date        time.Time `json:"date" bson:"date"`
	Balance     float64   `json:"balance" bson:"balance"`
}

type GetMonthlyReportReq struct {
	PropertyID string     `param:"propertyID" validate:"required"`
	Month      time.Month `query:"month" validate:"required"`
	Year       int        `query:"year" validate:"gte=1970,lte=2030"`
	Offset     int        `query:"offset" validate:"omitempty,gt=0"`
	Limit      int        `query:"limit" validate:"omitempty,gt=0"`
	NextToken  string     `query:"next_token"`
}
type GetMonthlyReportRes struct {
	StartingBalance float64               `json:"starting_balance"`
	Events          []*MonthlyReportEvent `json:"events"`
	NextToken       string                `json:"next_token"`
}

func (h *RestHandler) GetMonthlyReport(c echo.Context) error {
	req := &GetMonthlyReportReq{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if req.NextToken != "" {
		token, err := rest.DecodeNextToken(req.NextToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		req.Limit = token.Limit
		req.Offset = token.Offset
	}

	events, startingBalance, err := h.PropertyHandler.GetMonthlyReport(context.Background(), req.PropertyID, req.Month, req.Year, req.Offset, req.Limit)
	if err != nil {
		return err
	}

	mappedEvents := lo.Map(events, func(e *property.Event, _ int) *MonthlyReportEvent {
		return &MonthlyReportEvent{
			PropertyID:  e.PropertyID,
			EventAmount: e.EventAmount,
			Date:        e.Date,
			Balance:     e.PostEventBalance,
		}
	})
	res := &GetMonthlyReportRes{
		Events:          mappedEvents,
		StartingBalance: startingBalance,
	}

	if len(events) >= req.Limit {
		nextToken, err := rest.CreateNextToken(req.Limit, req.Offset+req.Limit)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		res.NextToken = nextToken
	}

	return c.JSON(200, res)
}
