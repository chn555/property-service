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

type GetEventsReq struct {
	PropertyID string    `param:"propertyID" validate:"required"`
	DateFrom   time.Time `query:"date_from"`
	DateTo     time.Time `query:"date_to"`
	SortOrder  string    `query:"sort_order" validate:"omitempty,oneof=asc desc"`
	AmountType string    `query:"amount_type" validate:"omitempty,oneof=expense income"`
	Offset     int       `query:"offset" validate:"omitempty,gt=0"`
	Limit      int       `query:"limit" validate:"omitempty,gt=0"`
	NextToken  string    `query:"next_token"`
}

type GetEventsRes struct {
	Events    []*Event `json:"events"`
	NextToken string   `json:"next_token"`
}

type Event struct {
	PropertyID  string    `json:"property_id,omitempty" bson:"property_id"`
	EventAmount float64   `json:"event_amount" bson:"event_amount"`
	Date        time.Time `json:"date" bson:"date"`
}

func (h *RestHandler) GetEvents(c echo.Context) error {
	req := &GetEventsReq{}
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
	sortOrder := property.Descending
	if req.SortOrder == "asc" {
		sortOrder = property.Ascending
	}

	amountType := property.All
	if req.AmountType == "income" {
		amountType = property.Income
	} else if req.AmountType == "expense" {
		amountType = property.Expense
	}

	events, err := h.PropertyHandler.GetPropertyEvents(context.Background(), req.PropertyID, req.DateFrom, req.DateTo, sortOrder, amountType, req.Offset, req.Limit)
	if err != nil {
		return err
	}

	mappedEvents := lo.Map(events, func(e *property.Event, _ int) *Event {
		return &Event{
			PropertyID:  e.PropertyID,
			EventAmount: e.EventAmount,
			Date:        e.Date,
		}
	})
	res := &GetEventsRes{
		Events: mappedEvents,
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

type SaveEventReq struct {
	PropertyID string  `param:"propertyID" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
}

func (h *RestHandler) SaveEvent(c echo.Context) error {
	req := &SaveEventReq{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	balance, err := h.PropertyHandler.SaveEvent(context.Background(), req.PropertyID, req.Amount, time.Now())
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]interface{}{"balance": balance})
}
