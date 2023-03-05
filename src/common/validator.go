package common

import (
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

	tv := &TranslatedValidator{
		validate:   validate,
		translator: translator,
	}

	if err := tv.registerValidation(
		"gender",
		"{0} is invalid gender",
		validateGender,
	); err != nil {
		return nil, err
	}

	if err := tv.registerValidation(
		"live-status",
		"{0} is invalid live status",
		validateLiveStatus,
	); err != nil {
		return nil, err
	}

	return tv, nil
}

func (tv *TranslatedValidator) ValidateStruct(r any) error {
	err := tv.validate.Struct(r)
	if err == nil {
		return nil
	}

	var serverError *validator.InvalidValidationError
	if errors.As(err, &serverError) {
		return errors.Wrap(serverError, "cannot validate struct")
	}

	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return errors.Wrap(err, "unexpected validation error type")
	}

	invalidParams := make([]InvalidRequestParameter, len(errs))
	for i, err := range errs {
		invalidParams[i] = InvalidRequestParameter{
			Name:   err.Field(),
			Reason: err.Translate(tv.translator),
		}
	}

	return NewValidationError(invalidParams...)
}

func (tv *TranslatedValidator) ValidateVar(key string, value any, tag string) error {
	err := tv.validate.Var(value, tag)
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
			Reason: strings.TrimSpace(err.Translate(tv.translator)),
		}
	}

	return NewValidationError(invalidParams...)
}

func (tv *TranslatedValidator) registerValidation(
	tag string,
	msg string,
	validatorFunc validator.Func,
) error {
	if err := tv.validate.RegisterValidation(tag, validatorFunc); err != nil {
		return err
	}

	err := tv.setErrorMessage(tag, msg)

	return err
}

func (tv *TranslatedValidator) setErrorMessage(
	tag string,
	msg string,
) error {
	return tv.validate.RegisterTranslation(
		tag,
		tv.translator,
		func(ut ut.Translator) error {
			return ut.Add(tag, msg, false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(tag, fe.Field())
			return t
		},
	)
}

var (
	Genders = []string{
		"MALE",
		"FEMALE",
		"OTHER",
	}
	LiveStatuses = []string{
		"ALIVE",
		"DEAD",
	}
)

func validateGender(fieldLevel validator.FieldLevel) bool {
	targetToCheck := fieldLevel.Field().String()

	if targetToCheck == "" {
		return true
	}

	for _, gender := range Genders {
		if gender == targetToCheck {
			return true
		}
	}

	return false
}

func validateLiveStatus(fieldLevel validator.FieldLevel) bool {
	targetToCheck := fieldLevel.Field().String()

	if targetToCheck == "" {
		return true
	}

	for _, status := range LiveStatuses {
		if status == targetToCheck {
			return true
		}
	}

	return false
}
