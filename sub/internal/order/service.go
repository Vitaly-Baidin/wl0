package order

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Service struct {
	repository *repository
}

func NewOrderService(database *pgxpool.Pool) *Service {
	repository := repository{database}
	return &Service{
		repository: &repository,
	}
}

func (s *Service) AddOrder(order Order) {
	s.repository.AddOrder(context.Background(), order)
}

func (s *Service) GetAllOrders() []Order {
	return s.repository.GetAllOrders(context.Background())
}

func (s *Service) GetOrderByUID(uid string) *Order {
	return s.repository.GetOrderByUID(context.Background(), uid)
}
