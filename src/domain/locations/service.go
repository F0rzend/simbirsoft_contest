package locations

import (
	"context"

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

func (s *Service) GetLocation(ctx context.Context, id int64) (*Entity, error) {
	return s.repository.GetLocation(ctx, id)
}
