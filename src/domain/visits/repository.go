package visits

import (
	"context"
	"time"

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

func (r *Repository) GetAnimalLocations(
	ctx context.Context,
	id int64,
	startDateTime *time.Time,
	endDateTime *time.Time,
	from int,
	size int,
) ([]*Entity, error) {
	return []*Entity{
		{
			ID:                           1,
			DateTimeOfVisitLocationPoint: time.Now(),
			LocationPointID:              1,
		},
	}, nil
}
