package service

import (
	"context"
	"github.com/Vitaly-Baidin/l0/pkg/config/viperconf"
	"github.com/Vitaly-Baidin/l0/sub/internal/domain"
	"github.com/Vitaly-Baidin/l0/sub/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/patrickmn/go-cache"
	"time"
)

type CacheService struct {
	repository *repository.CacheRepository
	Context    context.Context
	Cache      *cache.Cache
}

func NewCacheService(database *pgxpool.Pool, context context.Context, cfg *viperconf.Config) *CacheService {
	defaultExpiration := time.Duration(cfg.CacheConfig.DefaultExpiration) * time.Minute
	cleanupInterval := time.Duration(cfg.CacheConfig.CleanupInterval) * time.Minute

	c := cache.New(defaultExpiration, cleanupInterval)

	cacheRepository := repository.CacheRepository{Database: database}

	return &CacheService{
		repository: &cacheRepository,
		Context:    context,
		Cache:      c,
	}
}

func (s *CacheService) GetAllCacheFromDB() ([]domain.Cache, error) {
	return s.repository.GetAllCache(s.Context)
}

func (s *CacheService) SaveCache(key string, value any) error {
	s.Cache.Set(key, value, cache.DefaultExpiration)

	expiration := time.Duration(s.Context.Value("cache.expiration").(int))

	c := domain.Cache{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}

	return s.repository.AddCache(s.Context, c)
}

func (s *CacheService) RemoveCacheFromDB(key string) error {
	return s.repository.RemoveCache(s.Context, key)
}
