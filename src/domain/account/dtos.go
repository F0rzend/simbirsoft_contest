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

func (rr *RegistrationRequest) Bind(r *http.Request) error {
	tv, err := common.TranslatedValidatorFromRequest(r)
	if err != nil {
		return err
	}

	rr.FirstName = strings.TrimSpace(rr.FirstName)
	rr.LastName = strings.TrimSpace(rr.LastName)

	return tv.ValidateStruct(rr)
}

type Response struct {
	ID        uint   `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (r *Response) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type SearchParameters struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"emailName"`
	From      int    `json:"from" validate:"gte=0"`
	Size      int    `json:"size" validate:"gt=0"`
}

func NewSearchParameters(r *http.Request) (*SearchParameters, error) {
	tv, err := common.TranslatedValidatorFromRequest(r)
	if err != nil {
		return nil, err
	}

	values := r.URL.Query()

	var invalid []common.InvalidRequestParameter

	firstName := values.Get("firstName")
	lastName := values.Get("lastName")
	email := values.Get("email")

	from, invalidParameter := common.GetIntFromQuery(values, "from", 0)
	if invalidParameter != nil {
		invalid = append(invalid, *invalidParameter)
	}

	size, invalidParameter := common.GetIntFromQuery(values, "from", 10)
	if invalidParameter != nil {
		invalid = append(invalid, *invalidParameter)
	}

	if len(invalid) != 0 {
		return nil, common.NewValidationError(invalid...)
	}

	params := &SearchParameters{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		From:      from,
		Size:      size,
	}

	if err := tv.ValidateStruct(params); err != nil {
		return nil, err
	}

	return params, nil
}

type ResponseList []*Response

func (ResponseList) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
