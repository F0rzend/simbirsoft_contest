package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

func Bind(r *http.Request, body render.Binder) error {
	err := render.Bind(r, body)
	if err == nil {
		return nil
	}

	var jsonUnmarshalError *json.UnmarshalTypeError
	if errors.As(err, &jsonUnmarshalError) {
		return NewValidationError(InvalidRequestParameter{
			Name: jsonUnmarshalError.Field,
			Reason: fmt.Sprintf(
				"The field value must be %q type",
				jsonUnmarshalError.Value,
			),
		})
	}

	if errors.Is(err, io.EOF) {
		return &HTTPError{
			Type:   RequestBodyExpectedErrorType,
			Status: http.StatusBadRequest,
			Title:  "Request body expected.",
		}
	}

	return err
}

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

type Checker func(ctx context.Context, login string, password string) error

func BasicAuth(checkFn Checker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			login, password, ok := r.BasicAuth()
			if !ok {
				RenderError(w, r, NewUnauthorizedError("Invalid header format"))
				return
			}

			if err := checkFn(r.Context(), login, password); err != nil {
				RenderError(w, r, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
