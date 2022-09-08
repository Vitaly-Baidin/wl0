package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Vitaly-Baidin/l0/sub/internal/domain"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CacheRepository struct {
	Database *pgxpool.Pool
}

func (r *CacheRepository) AddCache(ctx context.Context, cache domain.Cache) error {
	_, err := r.Database.Exec(ctx, addCacheQuery,
		cache.Key, cache.Value, cache.Expiration)
	if err != nil {
		return fmt.Errorf("failed to add item cache to db: %v\n", err)
	}
	return nil
}

func (r *CacheRepository) GetAllCache(ctx context.Context) ([]domain.Cache, error) {
	var caches []domain.Cache
	rows, err := r.Database.Query(ctx, getAllCacheQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to find all cache from db: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		cache := domain.Cache{}
		err = rows.Scan(&cache.Key, &cache.Value, &cache.Expiration)
		if err != nil {
			return nil, fmt.Errorf("failed to convert from db: %v\n", err)
		}
		cache.Value, err = convertJsonToOrder(cache.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert from db: %v\n", err)
		}

		caches = append(caches, cache)
	}
	if rows.Err() != nil {
		return nil, pgx.ErrNoRows
	}
	return caches, nil
}

func (r *CacheRepository) RemoveCache(ctx context.Context, key string) error {
	_, err := r.Database.Exec(ctx, deleteCacheByKeyQuery, key)
	if err != nil {
		return fmt.Errorf("failed to delete cache in db: %v\n", err)
	}
	return nil
}

func convertJsonToOrder(value any) (*domain.Order, error) {
	orderType, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var o domain.Order
	err = json.Unmarshal(orderType, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

const (
	addCacheQuery = `INSERT INTO caches (Key, Value, Expiration)
					 VALUES ($1, $2, $3)
					 ON CONFLICT (Key)
					 DO UPDATE SET value = Excluded.Value, expiration = Excluded.Expiration;`
	getAllCacheQuery      = `SELECT Key, Value, Expiration FROM caches;`
	deleteCacheByKeyQuery = "DELETE FROM caches WHERE Key = $1;"
)
