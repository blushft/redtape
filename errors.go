package redtape

import (
	"net/http"

	"github.com/pkg/errors"
)

// Error is a customized error implementation with additional context for policy evaluation
type Error struct {
	code   int
	id     string
	reason string
	status string
	error
}

// StatusCode can contain application or standard integer codes eg http 401
func (e *Error) StatusCode() int {
	return e.code
}

// RequestID allows errors to be tracked against custom request ids
func (e *Error) RequestID() string {
	return e.id
}

// Status is a text explanation of the error code
func (e *Error) Status() string {
	return e.status
}

// Reason contains information about the policy decision that resulted in the error
func (e *Error) Reason() string {
	return e.reason
}

// NewErrRequestDeniedExplicit returns an error with for explicit denials
func NewErrRequestDeniedExplicit(err error) error {
	if err == nil {
		err = errors.New("request denied")
	}

	return errors.WithStack(&Error{
		error:  err,
		code:   http.StatusForbidden,
		status: http.StatusText(http.StatusForbidden),
		reason: "request denied because a policy explicitly forbids it",
	})
}

// NewErrRequestDeniedImplicit returns an error with for implicit denials (no policy)
func NewErrRequestDeniedImplicit(err error) error {
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
