package account

import (
	"net/http"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type RegistrationRequest struct {
	FirstName string `json:"firstName" validate:"required,alphanum"`
	LastName  string `json:"lastName" validate:"required,alphanum"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,alphanum"`
}

func (rr *RegistrationRequest) Bind(request *http.Request) error {
	validator, err := common.TranslatedValidatorFromRequest(request)
	if err != nil {
		return err
	}

	return validator.ValidateStruct(rr)
}

type Response struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (rr *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
