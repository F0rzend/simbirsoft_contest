// Package common
// Module errors
// RFC-7807 error handling
package common

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/hlog"
)

const (
	InternalServerErrorType = "InternalServerError"
	ValidationErrorType     = "ValidationError"
	ConflictErrorType       = "ConflictError"
	NotFoundErrorType       = "NotFoundError"
)

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	//nolint:errorlint
	if renderer, ok := err.(render.Renderer); ok {
		if renderError := render.Render(w, r, renderer); renderError != nil {
			RenderError(w, r, NewInternalServerError(renderError))
		}
		return
	}

	jsonUnmarshalError := new(json.UnmarshalTypeError)
	if errors.As(err, &jsonUnmarshalError) {
		_, ok := err.(*json.UnmarshalTypeError)
		if !ok {
			RenderError(w, r, NewInternalServerError(errors.New("json.UnmarshalTypeError was expected")))
		}
	}

	RenderError(w, r, NewInternalServerError(err))
}

type BaseHTTPError struct {
	ErrorType string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
}

func (e *BaseHTTPError) Error() string {
	return e.ErrorType
}

func (e *BaseHTTPError) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Status)

	return nil
}

func NewBaseHTTPError(errorType, title string, status int) *BaseHTTPError {
	return &BaseHTTPError{
		ErrorType: errorType,
		Title:     title,
		Status:    status,
	}
}

type InternalServerError struct {
	*BaseHTTPError

	err error
}

func NewInternalServerError(err error) *InternalServerError {
	return &InternalServerError{
		BaseHTTPError: NewBaseHTTPError(
			InternalServerErrorType,
			"Error on our side.",
			http.StatusInternalServerError,
		),
		err: err,
	}
}

func (e *InternalServerError) Render(w http.ResponseWriter, r *http.Request) error {
	if err := e.BaseHTTPError.Render(w, r); err != nil {
		return err
	}

	logger := hlog.FromRequest(r)

	logger.Error().Err(e.err).Send()

	return nil
}

type ValidationError struct {
	*BaseHTTPError

	InvalidParams []InvalidRequestParameter `json:"invalidParams,omitempty"`
}

type InvalidRequestParameter struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func NewValidationError(invalid ...InvalidRequestParameter) *ValidationError {
	return &ValidationError{
		BaseHTTPError: NewBaseHTTPError(
			ValidationErrorType,
			"Your request parameters didn't validate.",
			http.StatusBadRequest,
		),
		InvalidParams: invalid,
	}
}

type ConflictError struct {
	*BaseHTTPError

	Description string `json:"description"`
}

func NewConflictError(description string) *ConflictError {
	return &ConflictError{
		BaseHTTPError: NewBaseHTTPError(
			ConflictErrorType,
			"A data conflict has occurred.",
			http.StatusConflict,
		),
		Description: description,
	}
}

type NotFoundError struct {
	*BaseHTTPError

	Description string `json:"description"`
}

func NewNotFoundError(title, description string) *NotFoundError {
	return &NotFoundError{
		BaseHTTPError: NewBaseHTTPError(
			NotFoundErrorType,
			title,
			http.StatusNotFound,
		),
		Description: description,
	}
}
