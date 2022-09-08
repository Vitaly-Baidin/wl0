package cache

import (
	"context"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repository struct {
	Database *pgxpool.Pool
}

func (r *repository) AddCache(ctx context.Context, cache Cache) {
	_, err := r.Database.Exec(ctx, addCacheQuery,
		cache.Key, cache.Value, cache.Expiration)
	if err != nil {
		zaplog.Logger.Errorf("failed to create cache by: %v\n", err)
	}
}

func (r *repository) GetAllCache(ctx context.Context) []Cache {
	var caches []Cache
	rows, err := r.Database.Query(ctx, getAllCacheQuery)
	if err == pgx.ErrNoRows {
		zaplog.Logger.Info("No rows")
		return []Cache{}
	} else if err != nil {
		zaplog.Logger.Errorf("failed to all found cache: %v\n", err)
		return []Cache{}
	}
	defer rows.Close()
	for rows.Next() {
		cache := Cache{}
		rows.Scan(&cache.Key, &cache.Value, &cache.Expiration)
		caches = append(caches, cache)
	}

	if rows.Err() != nil {
		zaplog.Logger.Errorf("failed to find all caches: %v\n", err)
		return []Cache{}
	}

	return caches
}

func (r *repository) removeCache(ctx context.Context, key string) {
	_, err := r.Database.Exec(ctx, deleteCacheByKeyQuery, key)
	if err != nil {
		zaplog.Logger.Errorf("failed to delete cache: %v\n", err)
	}
}

const (
	addCacheQuery = `INSERT INTO caches (Key, Value, Expiration)
					 VALUES ($1, $2, $3)`
	getAllCacheQuery      = `SELECT Key, Value, Expiration FROM caches;`
	deleteCacheByKeyQuery = "DELETE FROM caches WHERE Key = $1;"
)
