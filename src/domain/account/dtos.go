package account

import (
	"net/http"
	"strings"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type RegistrationRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,alphanum"`
}

func (rr *RegistrationRequest) Bind(request *http.Request) error {
	validator, err := common.TranslatedValidatorFromRequest(request)
	if err != nil {
		return err
	}

	rr.FirstName = strings.TrimSpace(rr.FirstName)
	rr.LastName = strings.TrimSpace(rr.LastName)

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
