package common

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	et "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	envalidator "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
)

type TranslatedValidator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func NewTranslatedValidator() (*TranslatedValidator, error) {
	validate := validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	en := et.New()
	uni := ut.New(en, en)

	translator, _ := uni.GetTranslator("en")
	if err := envalidator.RegisterDefaultTranslations(validate, translator); err != nil {
		return nil, errors.Wrap(err, "error while register validator translation")
	}

	return &TranslatedValidator{
		validate:   validate,
		translator: translator,
	}, nil
}

func (v *TranslatedValidator) ValidateStruct(r any) error {
	err := v.validate.Struct(r)
	if err == nil {
		return nil
	}

	var serverError *validator.InvalidValidationError
	if errors.As(err, &serverError) {
		return err
	}

	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return err
	}

	invalidParams := make([]InvalidRequestParameter, len(errs))
	for i, err := range errs {
		invalidParams[i] = InvalidRequestParameter{
			Name:   err.Field(),
			Reason: err.Translate(v.translator),
		}
	}

	return NewValidationError(invalidParams...)
}

func (v *TranslatedValidator) ValidateVar(key string, value any, tag string) error {
	err := v.validate.Var(value, tag)
	if err == nil {
		return nil
	}

	var serverError *validator.InvalidValidationError
	if errors.As(err, &serverError) {
		return err
	}

	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return err
	}

	invalidParams := make([]InvalidRequestParameter, len(errs))
	for i, err := range errs {
		invalidParams[i] = InvalidRequestParameter{
			Name:   key,
			Reason: strings.TrimSpace(err.Translate(v.translator)),
		}
	}

	return NewValidationError(invalidParams...)
}

type ctxKey struct {
	id string
}

var translatedValidatorCtxKey = &ctxKey{"translatedValidator"}

type Middleware func(handler http.Handler) http.Handler

func TranslatedValidatorCtxMiddleware(translatedValidator *TranslatedValidator) Middleware {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), translatedValidatorCtxKey, translatedValidator))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func TranslatedValidatorFromRequest(r *http.Request) (*TranslatedValidator, error) {
	translatedValidator, ok := r.Context().Value(translatedValidatorCtxKey).(*TranslatedValidator)
	if !ok {
		return nil, errors.New("Cannot get validator from request context")
	}

	return translatedValidator, nil
}
