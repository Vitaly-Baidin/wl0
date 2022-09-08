package cache

import (
	"github.com/Vitaly-Baidin/l0/sub/internal/order"
	"time"
)

type Cache struct {
	ID         int
	Key        string
	Value      order.Order
	Expiration time.Duration
}
