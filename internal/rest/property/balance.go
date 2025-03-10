package property

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *RestHandler) getBalance(c echo.Context) error {
	propertyID := c.Param("propertyID")
	if propertyID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "property ID is required")
	}
	balance, err := h.PropertyHandler.GetBalance(context.Background(), propertyID)
	if err != nil {
		return err
	}
	return c.JSON(200, map[string]interface{}{"balance": balance})
}
