package types

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

func (s *Service) GetAnimalType(ctx context.Context, animalTypeID int64) (*Entity, error) {
	return s.repository.GetAnimalType(ctx, animalTypeID)
}
