package domain

import (
	"time"
)

type Cache struct {
	ID         int
	Key        string
	Value      any
	Expiration time.Duration
}
