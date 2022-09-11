package service

import (
	"context"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/entity"
	"github.com/Vitaly-Baidin/l0/pkg/repository"
	"github.com/Vitaly-Baidin/l0/pkg/util"
	"github.com/jackc/pgx/v4"
)

type OrderService struct {
	repository   repository.OrderRepository
	cacheService *CacheService
}

func NewOrderService(orderRepo repository.OrderRepository, cacheService *CacheService) *OrderService {
	return &OrderService{
		repository:   orderRepo,
		cacheService: cacheService,
	}
}

func (s *OrderService) AddOrder(ctx context.Context, order *entity.Order) error {
	err := s.repository.SaveOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed create order: %v", err)
	}

	cache := entity.Cache{
		Key:   *order.OrderUID,
		Value: order,
	}

	err = s.cacheService.SaveCache(ctx, cache)
	if err != nil {
		return fmt.Errorf("failed save to cache: %v", err)
	}

	return nil
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	return s.repository.GetAllOrders(ctx)
}

func (s *OrderService) GetOrderByUID(ctx context.Context, uid string) (*entity.Order, error) {
	value, found := s.cacheService.cache.Get(uid)
	if found {
		order, err := util.ConvertJsonToOrder(value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert cache from db: %v", err)
		}
		return order, nil
	}

	order, err := s.repository.GetOrderByUID(ctx, uid)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("order not found: %v", err)
	} else if err != nil {
		return nil, fmt.Errorf("failed order found: %v", err)
	}

	cache := entity.Cache{
		Key:   *order.OrderUID,
		Value: order,
	}

	err = s.cacheService.SaveCache(ctx, cache)
	if err != nil {
		return nil, fmt.Errorf("failed save cache: %v", err)
	}
	return order, nil
}
