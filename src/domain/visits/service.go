package visits

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

func (s *Service) GetAnimalLocations(
	ctx context.Context,
	id int64,
	startDateTime *time.Time,
	endDateTime *time.Time,
	from int,
	size int,
) ([]*Entity, error) {
	return s.repository.GetAnimalLocations(
		ctx,
		id,
		startDateTime,
		endDateTime,
		from,
		size,
	)
}
