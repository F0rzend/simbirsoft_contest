package account

type Service struct {
	repository *Repository
}

func NewService() *Service {
	return &Service{
		repository: NewRepository(),
	}
}

func (s *Service) Register() {}
