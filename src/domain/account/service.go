package account

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Service struct {
	repository *Repository
}

func NewService(config *common.Config) (*Service, error) {
	pool, err := pgxpool.New(context.Background(), config.Database.URI())
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialise postgres connection pool")
	}

	repository, err := NewRepository(pool)
	if err != nil {
		return nil, err
	}

	return &Service{
		repository: repository,
	}, nil
}

func (s *Service) Register(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
	rawPassword string,
) (*Entity, error) {
	isFreeEmail, err := s.repository.IsFreeEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !isFreeEmail {
		return nil, common.NewConflictError(fmt.Sprintf("Account with id %q already exists", email))
	}

	accountID, err := s.repository.NextAccountID(ctx)
	if err != nil {
		return nil, err
	}

	entity, err := NewEntity(accountID, firstName, lastName, email)
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "an error occurred while hashing the password")
	}

	if err := s.repository.Save(ctx, entity, string(passwordHash)); err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *Service) GetAccount(ctx context.Context, id uint) (*Entity, error) {
	return s.repository.GetAccount(ctx, id)
}

func (s *Service) Search(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
	from int,
	size int,
) ([]Entity, error) {
	return s.repository.Search(
		ctx,
		firstName,
		lastName,
		email,
		from,
		size,
	)
}

func (s *Service) Auth(
	ctx context.Context,
	email string,
	password string,
) error {
	hash, err := s.repository.GetPasswordHash(ctx, email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return common.NewUnauthorizedError("Wrong login or password.")
	}

	return nil
}
