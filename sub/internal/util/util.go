package util

import (
	"encoding/json"
	"github.com/Vitaly-Baidin/l0/sub/internal/domain"
)

func ConvertJsonToOrder(value any) (*domain.Order, error) {
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
