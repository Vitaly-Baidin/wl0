package service

import (
	"context"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/Vitaly-Baidin/l0/sub/internal/domain"
	"github.com/Vitaly-Baidin/l0/sub/internal/repository"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderService struct {
	repository   *repository.OrderRepository
	Context      context.Context
	cacheService *CacheService
}

func NewOrderService(database *pgxpool.Pool, context context.Context, cacheService *CacheService) *OrderService {
	orderRepository := repository.OrderRepository{Database: database}
	return &OrderService{
		repository:   &orderRepository,
		Context:      context,
		cacheService: cacheService,
	}
}

func (s *OrderService) AddOrder(order domain.Order) error {
	err := s.repository.AddOrder(s.Context, order)
	if err != nil {
		return fmt.Errorf("failed create order: %v\n", err)
	}
	return nil
}

func (s *OrderService) GetAllOrders() ([]domain.Order, error) {
	orders, err := s.repository.GetAllOrders(s.Context)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("not found: %v\n", err)
	} else if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *OrderService) GetOrderByUID(uid string) (*domain.Order, error) {
	value, found := s.cacheService.Cache.Get(uid)
	if found {
		zaplog.Logger.Infof("get value uid(%s) from cache", uid)
		return value.(*domain.Order), nil
	}
	order, err := s.repository.GetOrderByUID(s.Context, uid)
	if err == pgx.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}

	err = s.cacheService.SaveCache(uid, order)
	if err != nil {
		return nil, err
	}
	zaplog.Logger.Infof("get value uid(%s) from db", uid)
	return order, nil
}
