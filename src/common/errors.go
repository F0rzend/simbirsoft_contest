// Package common
// Module errors
// RFC-7807 error handling
package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/hlog"
)

// TODO: Replace to links with documentation
const (
	RequestBodyExpectedErrorType = "RequestBodyExpectedError"
	InternalServerErrorType      = "InternalServerError"
	ValidationErrorType          = "ValidationError"
	ConflictErrorType            = "ConflictError"
	NotFoundErrorType            = "NotFoundError"
	UnauthorizedErrorType        = "UnauthorizedError"
	ForbiddenErrorType           = "ForbiddenError"
)

//nolint:errname,stylecheck
var clientSideError = errors.New("client side error")

type HTTPError struct {
	Type     string `json:"type"`
	Status   int    `json:"status"`
	Title    string `json:"title,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`

	err error
}

func (e *HTTPError) Error() string {
	err := fmt.Sprintf("%s#%d", e.Type, e.Status)

	if e.Title != "" {
		err = fmt.Sprintf("%s: %s", err, e.Title)
	}

	return err
}

func (e *HTTPError) Unwrap() error {
	if e.err != nil {
		return e.err
	}

	return clientSideError
}

func (e *HTTPError) Render(_ http.ResponseWriter, r *http.Request) error {
	if e.err != nil {
		hlog.FromRequest(r).Error().Err(e.err).Send()
	}
	render.Status(r, e.Status)

	return nil
}

func NewInternalServerError(err error) error {
	return &HTTPError{
		Type:   InternalServerErrorType,
		Status: http.StatusInternalServerError,
		Title:  "Error on our side.",
		err:    err,
	}
}

func NewValidationError(invalid ...InvalidRequestParameter) error {
	return &struct {
		*HTTPError
		InvalidParams []InvalidRequestParameter `json:"invalidParams,omitempty"`
	}{
		HTTPError: &HTTPError{
			Type:   ValidationErrorType,
			Status: http.StatusBadRequest,
			Title:  "You request params didn't validate.",
		},
		InvalidParams: invalid,
	}
}

type InvalidRequestParameter struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func NewConflictError(detail string) error {
	return &HTTPError{
		Type:   ConflictErrorType,
		Status: http.StatusConflict,
		Title:  "A data conflict has occurred.",
		Detail: detail,
	}
}

func NewNotFoundError(title, detail string) error {
	return &HTTPError{
		Type:   NotFoundErrorType,
		Status: http.StatusNotFound,
		Title:  title,
		Detail: detail,
	}
}

func NewAuthorizationError(detail string) error {
	return &HTTPError{
		Type:   UnauthorizedErrorType,
		Status: http.StatusUnauthorized,
		Title:  "You should pass correct credentials in \"Authorization\" header.",
		Detail: detail,
	}
}

func NewForbiddenError(detail string) error {
	return &HTTPError{
		Type:   ForbiddenErrorType,
		Status: http.StatusForbidden,
		Title:  "Resource access denied.",
		Detail: detail,
	}
}
