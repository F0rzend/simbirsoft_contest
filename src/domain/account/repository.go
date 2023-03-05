package account

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(di *common.DependencyInjectionContainer) *Repository {
	return &Repository{
		db: di.Pool,
	}
}

func (r *Repository) NextAccountID(ctx context.Context) (uint, error) {
	query := `
		SELECT nextval('account_id')
	`

	var id uint
	if err := r.db.QueryRow(ctx, query).Scan(&id); err != nil {
		return 0, errors.Wrap(err, "an error occurred while getting next account id")
	}

	return id, nil
}

func (r *Repository) IsFreeEmail(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT NOT EXISTS(
		    SELECT 1 FROM account
		    WHERE email = @email
		)
	`

	var result bool
	if err := r.db.QueryRow(ctx, query, pgx.NamedArgs{
		"email": email,
	}).Scan(&result); err != nil {
		return false, errors.Wrap(err, "cannot check email freeness")
	}

	return result, nil
}

func (r *Repository) Save(
	ctx context.Context,
	entity *Entity,
	passwordHash string,
) error {
	query := `
		INSERT INTO account (id, first_name, last_name, email, password)
		VALUES (@id, @first_name, @last_name, @email, @password)
		ON CONFLICT ON CONSTRAINT account_pkey DO UPDATE
		SET first_name = excluded.first_name,
		    last_name = excluded.last_name,
		    email = excluded.email,
		    password = excluded.password
		WHERE account.id = excluded.id
	`

	if _, err := r.db.Exec(ctx, query, pgx.NamedArgs{
		"id":         entity.ID,
		"first_name": entity.FirstName,
		"last_name":  entity.LastName,
		"email":      entity.Email,
		"password":   passwordHash,
	}); err != nil {
		return errors.Wrap(err, "error while saving the account")
	}

	return nil
}

func (r *Repository) GetAccount(ctx context.Context, id uint) (*Entity, error) {
	query := `
		SELECT first_name, last_name, email FROM account
		WHERE id = @id
	`

	var firstName, lastName, email string
	if err := r.db.QueryRow(ctx, query, pgx.NamedArgs{
		"id": id,
	}).Scan(
		&firstName,
		&lastName,
		&email,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, common.NewNotFoundError(
				"Account not found.",
				fmt.Sprintf("Account with id %d not found.", id),
			)
		}

		return nil, errors.Wrap(err, "cannot get account by id")
	}

	entity, err := NewEntity(id, firstName, lastName, email)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r *Repository) Search(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
	from int,
	size int,
) ([]Entity, error) {
	query := `
		SELECT id, first_name, last_name, email
		FROM account
		WHERE (first_name ILIKE concat('%', $1::varchar, '%')) AND
			  (first_name ILIKE concat('%', $2::varchar, '%')) AND
			  (email ILIKE concat('%', $3::varchar, '%'))
		ORDER BY id
		LIMIT $4 OFFSET $5
	`

	var entities []Entity
	err := pgxscan.Select(
		ctx,
		r.db,
		&entities,
		query,
		firstName,
		lastName,
		email,
		size,
		from,
	)
	if err != nil {
		return nil, err
	}

	return entities, nil
}

//nolint:nonamedreturns
func (r *Repository) GetPasswordHash(ctx context.Context, email string) (id int, hash []byte, err error) {
	query := `
		SELECT id, password FROM account
		WHERE email = @email
	`

	if err = r.db.QueryRow(ctx, query, pgx.NamedArgs{
		"email": email,
	}).Scan(&id, &hash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil, common.NewNotFoundError(
				"Account not found.",
				fmt.Sprintf("Account with email %q not found.", email),
			)
		}

		return 0, nil, errors.Wrap(err, "cannot get account hash")
	}

	return id, hash, nil
}
