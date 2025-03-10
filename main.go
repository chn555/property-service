package main

import (
	"context"
	"github.com/chn555/property-service/pkg/db/mongo"
	"github.com/chn555/property-service/pkg/property"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	_ "github.com/swaggo/echo-swagger/example/docs"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {

	mongoClient, err := mongo.NewClient(context.TODO(), mongo.Config{
		URI:     "mongodb://localhost:27017",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		slog.Error("failed to connect to mongo db", err)
		os.Exit(1)
	}

	con := mongo.NewEventState(mongoClient, &mongo.EventStateConfig{
		DatabaseName:   "property",
		CollectionName: "events",
	})

	e := echo.New()
	e.Debug = true
	e.Validator = &CustomValidator{validator: validator.New()}
	h := property.NewHandler(con)
	e.POST("/property/:propertyID", func(c echo.Context) error {
		req := &SaveEventReq{}
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		balance, err := h.SaveEvent(context.Background(), req.PropertyID, req.Amount, time.Now())
		if err != nil {
			return err
		}

		return c.JSON(200, map[string]interface{}{"balance": balance})
	})

	e.GET("/property/:propertyID/events", func(c echo.Context) error {
		req := &GetEventsReq{}
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if req.NextToken != "" {
			token, err := decodeNextToken(req.NextToken)
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

		events, err := h.GetPropertyEvents(context.Background(), req.PropertyID, req.DateFrom, req.DateTo, sortOrder, amountType, req.Offset, req.Limit)
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
			nextToken, err := createNextToken(req.Limit, req.Offset+req.Limit)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			res.NextToken = nextToken
		}

		return c.JSON(200, res)
	})

	e.GET("/property/:propertyID/balance", func(c echo.Context) error {
		propertyID := c.Param("propertyID")
		if propertyID == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "property ID is required")
		}
		balance, err := h.GetBalance(context.Background(), propertyID)
		if err != nil {
			return err
		}
		return c.JSON(200, map[string]interface{}{"balance": balance})
	})
	if err := e.Start(":1323"); err != nil {
	}
	e.Logger.Fatal(e.Start(":1323"))
}

type SaveEventReq struct {
	PropertyID string  `param:"propertyID" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
	//Date       time.Time
}

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
