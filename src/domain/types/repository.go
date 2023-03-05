package types

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(di *common.DependencyInjectionContainer) *Repository {
	return &Repository{
		pool: di.Pool,
	}
}

func (r *Repository) GetAnimalType(ctx context.Context, animalTypeID int64) (*Entity, error) {
	return &Entity{
		ID:   animalTypeID,
		Type: "type",
	}, nil
}
