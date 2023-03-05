package animal

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

func (r *Repository) GetAnimal(_ context.Context, id int64) (*Entity, error) {
	now := time.Now()
	return &Entity{
		ID:                 id,
		Types:              []int64{1},
		Weight:             1,
		Length:             1,
		Height:             1,
		Gender:             "MALE",
		LifeStatus:         "ALIVE",
		ChippingDateTime:   &now,
		ChipperID:          1,
		ChippingLocationID: 1,
		VisitedLocations:   []int{1},
		DeathDateTime:      nil,
	}, nil
}

func (r *Repository) Search(
	ctx context.Context,
	startDateTime *time.Time,
	endDateTime *time.Time,
	chipperID int,
	chippingLocationID int64,
	lifeStatus string,
	gender string,
	from int,
	size int,
) ([]*Entity, error) {
	now := time.Now()
	return []*Entity{
		{
			ID:                 1,
			Types:              []int64{1},
			Weight:             1,
			Length:             1,
			Height:             1,
			Gender:             "MALE",
			LifeStatus:         "ALIVE",
			ChippingDateTime:   &now,
			ChipperID:          1,
			ChippingLocationID: 1,
			VisitedLocations:   []int{1},
			DeathDateTime:      nil,
		},
	}, nil
}
