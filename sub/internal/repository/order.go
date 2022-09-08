package repository

import (
	"context"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/Vitaly-Baidin/l0/sub/internal/domain"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrderRepository struct {
	Database *pgxpool.Pool
}

func (r *OrderRepository) Start() {
	fmt.Println("Project godb started!")
}

func (r *OrderRepository) AddOrder(ctx context.Context, order domain.Order) error {
	_, err := r.Database.Exec(
		ctx,
		addOrderQuery,
		order.OrderUID,
		order.TrackNumber, order.Entry, order.DeliveryData,
		order.PaymentData, order.ItemsData, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	rows, err := r.Database.Query(ctx, getAllOrdersQuery)
	if err != nil {
		zaplog.Logger.Errorf("failed to all found orders: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		order := domain.Order{}
		err = rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.DeliveryData,
			&order.PaymentData, &order.ItemsData, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) GetOrderByUID(ctx context.Context, uid string) (*domain.Order, error) {
	order := domain.Order{}

	err := r.Database.QueryRow(ctx, getOrderByUIDQuery, uid).
		Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.DeliveryData,
			&order.PaymentData, &order.ItemsData, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

const (
	addOrderQuery = `INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, 
                    					locale, internal_signature, customer_id, 
                    					delivery_service, shardkey, sm_id, date_created, 
                    					oof_shard)
					 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
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
