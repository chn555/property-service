package rest

import (
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

func NewServer(handlerRegisters ...func(e *echo.Echo) *echo.Echo) *echo.Echo {
	e := echo.New()
	e.Debug = true
	e.Validator = &CustomValidator{validator: validator.New()}

	lo.ForEach(handlerRegisters, func(handler func(e *echo.Echo) *echo.Echo, _ int) {
		handler(e)
	})
	return e
}
