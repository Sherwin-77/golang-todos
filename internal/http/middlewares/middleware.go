package middlewares

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Middleware struct {
}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) ValidateUUID(params []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, param := range params {
				if _, err := uuid.Parse(c.Param(param)); err != nil {
					return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
				}
			}

			return next(c)
		}
	}
}
