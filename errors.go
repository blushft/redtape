package redtape

import (
	"net/http"

	"github.com/pkg/errors"
)

type Error struct {
	code   int
	id     string
	reason string
	status string
	error
}

func (e *Error) StatusCode() int {
	return e.code
}

func (e *Error) RequestID() string {
	return e.id
}

func (e *Error) Status() string {
	return e.status
}

func (e *Error) Reason() string {
	return e.reason
}

func NewErrRequestDenied(err error) error {
	if err == nil {
		err = errors.New("request denied")
	}

	return errors.WithStack(&Error{
		error:  err,
		code:   http.StatusForbidden,
		status: http.StatusText(http.StatusForbidden),
		reason: "request denied because no matching policy was found",
	})
}
