// Package common
// Module errors
// RFC-7807 error handling
package common

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	RequestBodyExpectedErrorType = "RequestBodyExpectedError"
	InternalServerErrorType      = "InternalServerError"
	ValidationErrorType          = "ValidationError"
	ConflictErrorType            = "ConflictError"
	NotFoundErrorType            = "NotFoundError"
)

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	var httpError *HTTPError
	if errors.As(err, &httpError) {
		if renderError := render.Render(w, r, httpError); renderError != nil {
			RenderError(w, r, NewInternalServerError(renderError))
		}
		return
	}

	RenderError(w, r, NewInternalServerError(err))
}

type HTTPError struct {
	Type     string `json:"type"`
	Status   int    `json:"status"`
	Title    string `json:"title,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`

	err error
}

func (e *HTTPError) Error() string {
	return e.Title
}

func (e *HTTPError) Unwrap() error {
	if e.err != nil {
		return e.err
	}

	return e
}

func (e *HTTPError) Render(_ http.ResponseWriter, r *http.Request) error {
	if e.err != nil {
		zerolog.Ctx(r.Context()).Error().Err(e.err).Send()
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
