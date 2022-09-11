package repository

import (
	"context"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/entity"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	addCacheQuery = `INSERT INTO caches (Key, Value, Expiration)
					 VALUES ($1, $2, $3)
					 ON CONFLICT (Key)
					 DO UPDATE SET value = EXCLUDED.Value, expiration = EXCLUDED.Expiration;`
	getAllCacheQuery      = `SELECT Key, Value, Expiration FROM caches;`
	deleteCacheByKeyQuery = "DELETE FROM caches WHERE Key = $1;"
)

type CacheRepository interface {
	SaveCache(ctx context.Context, cache entity.Cache) error
	GetAllCaches(ctx context.Context) ([]entity.Cache, error)
	RemoveCache(ctx context.Context, key string) error
}

func NewCacheRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		database: pool,
	}
}

func (r *Repository) SaveCache(ctx context.Context, cache entity.Cache) error {
	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to create tx: %v", err)
	}

	var sqlArgs = []any{cache.Key, cache.Value, cache.Expiration}

	_, err = tx.Exec(ctx, addCacheQuery, sqlArgs...)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return fmt.Errorf("rollback err: %v, err: %v", rollbackErr, err)
		}

		return fmt.Errorf("failed to save cache to db: %v", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed commit tx: %v", err)
	}
	return nil
}

func (r *Repository) GetAllCaches(ctx context.Context) ([]entity.Cache, error) {
	var caches []entity.Cache

	rows, err := r.database.Query(ctx, getAllCacheQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to find all cache from db: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		cache := entity.Cache{}
		err = rows.Scan(&cache.Key, &cache.Value, &cache.Expiration)
		if err != nil {
			return nil, fmt.Errorf("failed scan cache from db: %v", err)
		}
		// TODO deprecated
		//order, err := util.ConvertJsonToOrder(cache.Value)
		//if err != nil {
		//	return nil, err
		//}
		//cache.Value = order
		//if err != nil {
		//	return nil, fmt.Errorf("failed to convert from db: %v\n", err)
		//}
		caches = append(caches, cache)
	}
	return caches, nil
}

func (r *Repository) RemoveCache(ctx context.Context, key string) error {
	_, err := r.database.Exec(ctx, deleteCacheByKeyQuery, key)
	if err != nil {
		return fmt.Errorf("failed to delete cache in db: %v\n", err)
	}
	return nil
}
