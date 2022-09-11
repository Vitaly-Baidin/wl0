package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Vitaly-Baidin/l0/config"
	"github.com/Vitaly-Baidin/l0/logging"
	"github.com/Vitaly-Baidin/l0/pkg/entity"
	"github.com/Vitaly-Baidin/l0/pkg/repository"
	"github.com/patrickmn/go-cache"
	"time"
)

type CacheService struct {
	repository repository.CacheRepository
	cache      *cache.Cache
}

func NewCacheService(cacheRepo repository.CacheRepository) *CacheService {
	return &CacheService{
		repository: cacheRepo,
	}
}

func (s *CacheService) InitCache(log logging.Logger, cfg config.Config) {
	defaultExpiration := time.Duration(cfg.CacheConfig.DefaultExpiration) * time.Minute
	cleanupInterval := time.Duration(cfg.CacheConfig.CleanupInterval) * time.Minute

	c := cache.New(defaultExpiration, cleanupInterval)
	c.OnEvicted(func(key string, v interface{}) {
		err := s.RemoveCache(context.TODO(), key)
		if err != nil {
			log.Errorf("failed initialization cache: %v", err)
		}
	})

	s.cache = c
}

func (s *CacheService) LoadCacheFromDB(ctx context.Context) error {
	if s.cache == nil {
		return errors.New("cache service not initialized")
	}
	caches, err := s.GetAllCaches(ctx)
	if err != nil {
		return fmt.Errorf("failed to load cache from db: %v", err)
	}

	for _, v := range caches {
		s.cache.Set(v.Key, v.Value, v.Expiration*time.Minute)
	}

	return nil
}

func (s *CacheService) SaveCache(ctx context.Context, c entity.Cache) error {
	s.cache.Set(c.Key, c.Value, cache.DefaultExpiration)
	return s.repository.SaveCache(ctx, c)
}

func (s *CacheService) GetAllCaches(ctx context.Context) ([]entity.Cache, error) {
	return s.repository.GetAllCaches(ctx)
}

func (s *CacheService) RemoveCache(ctx context.Context, key string) error {
	s.cache.Delete(key)
	return s.repository.RemoveCache(ctx, key)
}
