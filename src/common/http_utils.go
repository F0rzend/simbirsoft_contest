package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
	if errors.Is(err, clientSideError) {
		var renderer render.Renderer
		if !errors.As(err, &renderer) {
			RenderError(w, r, NewInternalServerError(errors.New("client side error must be renderer")))
		}

		if renderError := render.Render(w, r, renderer); renderError != nil {
			RenderError(w, r, NewInternalServerError(renderError))
		}
		return
	}

	RenderError(w, r, NewInternalServerError(err))
}

func GetIntFromRequest(r *http.Request, key string) (int, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return 0, NewValidationError(InvalidRequestParameter{
			Name:   key,
			Reason: fmt.Sprintf("%q is required", key),
		})
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, NewValidationError(InvalidRequestParameter{
			Name:   key,
			Reason: fmt.Sprintf("int expected, but got %s", param),
		})
	}

	return value, nil
}

func GetInt64FromRequest(r *http.Request, key string) (int64, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return 0, NewValidationError(InvalidRequestParameter{
			Name:   key,
			Reason: fmt.Sprintf("%q is required", key),
		})
	}

	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, NewValidationError(InvalidRequestParameter{
			Name:   key,
			Reason: fmt.Sprintf("int expected, but got %s", param),
		})
	}

	return value, nil
}

func GetIntFromQuery(
	values url.Values,
	key string,
	defaultValue int,
) (int, *InvalidRequestParameter) {
	param := values.Get(key)

	if param == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return 0, &InvalidRequestParameter{
			Name:   key,
			Reason: fmt.Sprintf("%s must be a valid number", key),
		}
	}

	return value, nil
}

func GetInt64FromQuery(
	values url.Values,
	key string,
	defaultValue int64,
) (int64, *InvalidRequestParameter) {
	param := values.Get(key)

	if param == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, &InvalidRequestParameter{
			Name:   key,
			Reason: fmt.Sprintf("%s must be a valid number", key),
		}
	}

	return value, nil
}

func GetDatetimeFromQuery(
	values url.Values,
	key string,
) (*time.Time, *InvalidRequestParameter) {
	dateTime := values.Get(key)
	if dateTime == "" {
		return nil, nil
	} else {
		dateTime, err := time.Parse(time.RFC3339, dateTime)
		if err != nil {
			return nil, &InvalidRequestParameter{
				Name:   key,
				Reason: "must be a datetime in RFC-3339 format",
			}
		}
		return &dateTime, nil
	}
}
