package cache

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/patrickmn/go-cache"
	"time"
)

type Service struct {
	repository *repository
	Cache      *cache.Cache
}

func NewCacheService(database *pgxpool.Pool, cache *cache.Cache) *Service {
	repository := repository{database}
	return &Service{
		repository: &repository,
		Cache:      cache,
	}
}

func (s *Service) GetAllCacheFromDB() []Cache {
	return s.repository.GetAllCache(context.Background())
}

func (s *Service) SaveCache(key string, value any, duration time.Duration) {
	s.Cache.Set(key, value, duration)
	c := Cache{
		Key:        key,
		Value:      value,
		Expiration: duration,
	}
	s.repository.AddCache(context.Background(), c)
}

func (s *Service) RemoveCacheFromDB(key string) {
	s.repository.removeCache(context.Background(), key)
}
