package domain

import (
	"time"
)

type Cache struct {
	Key        string
	Value      any
	Expiration time.Duration
}
