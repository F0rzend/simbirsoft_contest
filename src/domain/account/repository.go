package account

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) (*Repository, error) {
	if err := pool.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	return &Repository{
		db: pool,
	}, nil
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

	zerolog.Ctx(ctx).Debug().Bool("is_free", result).Str("email", email).Send()

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
		"email":      entity.Email.Address,
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
			return nil, ErrNotFound
		}

		return nil, errors.Wrap(err, "cannot get account by id")
	}

	entity, err := NewEntity(id, firstName, lastName, email)
	if err != nil {
		return nil, err
	}

	return entity, nil
}
