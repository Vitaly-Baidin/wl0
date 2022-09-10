package domain

import (
	"time"
)

type Order struct {
	OrderUID          *string   `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	DeliveryData      any       `json:"delivery,omitempty"`
	PaymentData       any       `json:"payment,omitempty"`
	ItemsData         []any     `json:"items,omitempty"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	Shardkey          string    `json:"shardkey"`
	SmID              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`
	OofShard          string    `json:"oof_shard"`
}
