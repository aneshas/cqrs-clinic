package errs

import (
	"errors"
	"github.com/aneshas/eventstore/aggregate"
	"github.com/labstack/echo/v4"
	"net/http"
)

type httpError struct {
	Code  int `json:"code"`
	Error any `json:"error"`
} // @name HTTPError

// ErrorHandler handles application errors
func ErrorHandler(err error, c echo.Context) {
	var appErr Error

	code := http.StatusInternalServerError

	if errors.As(err, &appErr) {
		code = http.StatusBadRequest

		_ = c.JSON(
			code,
			httpError{
				Code:  code,
				Error: appErr.Error(),
			},
		)

		return
	}

	if errors.Is(err, aggregate.ErrAggregateNotFound) {
		code = http.StatusNotFound

		_ = c.JSON(
			code,
			httpError{
				Code:  code,
				Error: "not found",
			},
		)
	}

	_ = c.JSON(
		code,
		httpError{
			Code:  code,
			Error: err.Error(),
		},
	)
}

// E creates new error
func E(err error) Error {
	return Error{
		err: err,
	}
}

// Error represents a domain error
type Error struct {
	err error
}

// Error returns error message
func (ve Error) Error() string {
	return ve.err.Error()
}
