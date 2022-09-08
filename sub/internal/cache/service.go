package cache

import (
	"context"
	"github.com/Vitaly-Baidin/l0/sub/internal/order"
	"time"
)

type Service struct {
	repository Repository
}

func NewCacheService(repository Repository) *Service {
	Repository{}
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetAllCacheFromDB() []Cache {
	return s.repository.GetAllCache(context.Background())
}

func (s *Service) SaveCacheToDB(key string, value order.Order, duration time.Duration) {
	cache := Cache{
		Key:        key,
		Value:      value,
		Expiration: duration,
	}
	s.repository.AddCache(context.Background(), cache)
}

func (s *Service) RemoveCacheFromDB(key string) {
	s.repository.removeCache(context.Background(), key)
}
