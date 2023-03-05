package animal

import (
	"context"
	"time"

	"github.com/F0rzend/simbirsoft_contest/src/common"
)

type Service struct {
	repository *Repository
}

func NewService(di *common.DependencyInjectionContainer) *Service {
	return &Service{
		repository: NewRepository(di),
	}
}

func (s *Service) GetAnimal(ctx context.Context, id int64) (*Entity, error) {
	return s.repository.GetAnimal(ctx, id)
}

func (s *Service) Search(
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
	return s.repository.Search(
		ctx,
		startDateTime,
		endDateTime,
		chipperID,
		chippingLocationID,
		lifeStatus,
		gender,
		from,
		size,
	)
}
