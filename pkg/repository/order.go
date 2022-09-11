package repository

import (
	"context"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/entity"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	addOrderQuery = `INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, 
                    					locale, internal_signature, customer_id, 
                    					delivery_service, shardkey, sm_id, date_created, 
                    					oof_shard)
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
					 ON CONFLICT (order_uid) DO UPDATE
					 SET track_number = excluded.track_number, entry = excluded.entry, delivery = excluded.delivery,
						 payment = excluded.payment, items = excluded.items, locale = excluded.locale, 
						 internal_signature = excluded.internal_signature, customer_id = excluded.customer_id,
						 delivery_service = excluded.delivery_service, shardkey = excluded.shardkey, 
						 sm_id = excluded.sm_id, date_created = excluded.date_created, oof_shard = excluded.oof_shard;`
	getAllOrdersQuery = `SELECT order_uid, track_number, entry, delivery, payment, items, 
                    					locale, internal_signature, customer_id, 
                    					delivery_service, shardkey, sm_id, date_created, 
                    					oof_shard 
						 FROM orders;`
	getOrderByUIDQuery = `SELECT order_uid, track_number, entry, delivery, payment, items, 
                    					locale, internal_signature, customer_id, 
                    					delivery_service, shardkey, sm_id, date_created, 
                    					oof_shard 
						  FROM orders
						  WHERE order_uid=$1
						  LIMIT 1;`
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *entity.Order) error
	GetAllOrders(ctx context.Context) ([]entity.Order, error)
	GetOrderByUID(ctx context.Context, uid string) (*entity.Order, error)
}

func NewOrderRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		database: pool,
	}
}

func (r *Repository) SaveOrder(ctx context.Context, order *entity.Order) error {
	tx, err := r.database.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to create tx: %v", err)
	}

	var sqlArgs = []any{order.OrderUID, order.TrackNumber, order.Entry, order.DeliveryData,
		order.PaymentData, order.ItemsData, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID,
		order.DateCreated, order.OofShard}

	_, err = tx.Exec(ctx, addOrderQuery, sqlArgs...)
	if err != nil {
		rollbackErr := tx.Rollback(ctx)
		if rollbackErr != nil {
			return fmt.Errorf("rollback err: %v, err: %v", rollbackErr, err)
		}

		return fmt.Errorf("failed to save order to db: %v", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed commit tx: %v", err)
	}
	return nil
}

func (r *Repository) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	var orders []entity.Order
	rows, err := r.database.Query(ctx, getAllOrdersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to find all orders from db: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		order := entity.Order{}
		err = rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.DeliveryData,
			&order.PaymentData, &order.ItemsData, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, fmt.Errorf("failed scan order from db: %v", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *Repository) GetOrderByUID(ctx context.Context, uid string) (*entity.Order, error) {
	o := entity.Order{}

	err := r.database.QueryRow(ctx, getOrderByUIDQuery, uid).
		Scan(&o.OrderUID, &o.TrackNumber, &o.Entry, &o.DeliveryData,
			&o.PaymentData, &o.ItemsData, &o.Locale,
			&o.InternalSignature, &o.CustomerID, &o.DeliveryService,
			&o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard)
	if err != nil {
		return nil, err
	}

	if o.OrderUID == nil {
		return nil, pgx.ErrNoRows
	}

	return &o, nil
}
