package property

import (
	"github.com/chn555/property-service/pkg/property"
	"github.com/labstack/echo/v4"
)

type RestHandler struct {
	PropertyHandler *property.Handler
}

func NewRestHandler(propertyHandler *property.Handler) *RestHandler {
	return &RestHandler{PropertyHandler: propertyHandler}
}

func (h *RestHandler) RegisterHandlers(e *echo.Echo) *echo.Echo {
	g := e.Group("/property")
	g.POST("/:propertyID", h.SaveEvent)
	g.GET("/:propertyID/events", h.GetEvents)
	g.GET("/:propertyID/monthly_report", h.GetMonthlyReport)
	g.GET("/:propertyID/balance", h.getBalance)

	return e
}
