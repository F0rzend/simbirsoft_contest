package account

import (
	"net/mail"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Entity struct {
	ID        uint
	FirstName string
	LastName  string
	Email     string
}

func NewEntity(id uint, firstName, lastName, email string) (*Entity, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, common.NewValidationError(common.InvalidRequestParameter{
			Name:   "email",
			Reason: "email is invalid",
		})
	}

	return &Entity{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}, nil
}
