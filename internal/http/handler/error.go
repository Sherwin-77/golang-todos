package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/go-echo-template/pkg/response"
)

func extractColumnNameFromDetail(detail string) string {
	start := strings.Index(detail, "(") + 1
	end := strings.Index(detail, ")")
	if start > 0 && end > start {
		return detail[start:end]
	}
	return "Field"
}

func HTTPErrorHandler(err error, c echo.Context) {
	var code int
	var message interface{}

	var pgerr *pgconn.PgError
	var he *echo.HTTPError
	var ve validator.ValidationErrors

	if errors.As(err, &he) {
		code = he.Code
		message = he.Message
	} else if errors.As(err, &pgerr) {
		code = http.StatusUnprocessableEntity
		if pgerr.Code == "23505" {
			column := extractColumnNameFromDetail(pgerr.Detail)
			message = fmt.Sprintf("%s already exists", column)
		} else {
			message = pgerr.Message
		}
	} else if errors.As(err, &ve) {
		code = http.StatusUnprocessableEntity
		fieldErr := ve[0]
		switch fieldErr.Tag() {
		case "required":
			message = fieldErr.Field() + " is required"
		case "email":
			message = fieldErr.Field() + " is not a valid email"
		case "gte":
			message = fieldErr.Field() + " must be greater than or equal to " + fieldErr.Param()
		case "lte":
			message = fieldErr.Field() + " must be less than or equal to " + fieldErr.Param()
		case "uuid":
			message = fieldErr.Field() + " is not a valid UUID"
		case "oneof":
			message = fieldErr.Field() + " must be one of " + fieldErr.Param()
		default:
			message = fieldErr.Field() + " is not valid"
		}
	} else {
		code = http.StatusInternalServerError
		message = err.Error()
		c.Logger().Error(err)
	}

	if !c.Response().Committed {
		_ = c.JSON(code, response.NewResponse(code, message.(string), nil, nil))
	}
}
