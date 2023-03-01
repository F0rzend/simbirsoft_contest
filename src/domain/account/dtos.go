package account

import (
	"net/http"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type RegistrationRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

func (r *RegistrationRequest) Bind(req *http.Request) error {
	validator, err := common.TranslatedValidatorFromRequest(req)
	if err != nil {
		return err
	}

	return validator.ValidateStruct(r)
}
