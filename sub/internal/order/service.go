package order

type Service struct {
	repository Repository
}

func NewOrderService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}
