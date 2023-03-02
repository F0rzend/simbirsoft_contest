package account

import (
	"net/mail"

	"github.com/pkg/errors"
)

type Entity struct {
	ID        uint
	FirstName string
	LastName  string
	Email     *mail.Address
}

func NewEntity(id uint, firstName, lastName, email string) (*Entity, error) {
	mailAddress, err := mail.ParseAddress(email)
	if err != nil {
		return nil, errors.Wrap(err, "Error parsing error")
	}

	return &Entity{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     mailAddress,
	}, nil
}
