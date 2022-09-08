package order

import (
	"context"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/logging/zaplog"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository struct {
	Database *pgxpool.Pool
}

func (r *Repository) Start() {
	fmt.Println("Project godb started!")
}

func (r *Repository) AddOrder(ctx context.Context, order Order) {
	_, err := r.Database.Exec(
		ctx,
		addOrderQuery,
		order.OrderUID,
		order.TrackNumber, order.Entry, order.DeliveryData,
		order.PaymentData, order.ItemsData, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		zaplog.Logger.Errorf("failed to create order by uid(%s): %v\n", order.OrderUID, err)
	}
}

func (r *Repository) GetAllOrders(ctx context.Context) []Order {
	var orders []Order
	rows, err := r.Database.Query(ctx, getAllOrdersQuery)
	if err == pgx.ErrNoRows {
		zaplog.Logger.Info("No rows")
		return []Order{}
	} else if err != nil {
		zaplog.Logger.Errorf("failed to all found orders: %v\n", err)
		return []Order{}
	}
	defer rows.Close()
	for rows.Next() {
		order := Order{}
		rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.DeliveryData,
			&order.PaymentData, &order.ItemsData, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
		orders = append(orders, order)
	}

	if rows.Err() != nil {
		zaplog.Logger.Errorf("failed to find all orders: %v\n", err)
		return []Order{}
	}

	return orders
}

func (r *Repository) GetOrderByUID(ctx context.Context, uid string) *Order {
	order := Order{}

	err := r.Database.QueryRow(ctx, getOrderByUIDQuery, uid).
		Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.DeliveryData,
			&order.PaymentData, &order.ItemsData, &order.Locale,
			&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
			&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard)
	if err != nil {
		zaplog.Logger.Errorf("failed to find order by uid(%s): %v\n", uid, err)
		return &Order{}
	}

	return &order
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
