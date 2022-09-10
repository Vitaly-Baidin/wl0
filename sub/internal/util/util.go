package util

import (
	"encoding/json"
	"fmt"
	"github.com/Vitaly-Baidin/l0/sub/internal/domain"
	"reflect"
)

func ConvertJsonToOrder(value any) (o *domain.Order, err error) {
	valueType := reflect.TypeOf(value).Kind()
	if valueType == reflect.String {
		valueString := fmt.Sprintf("%v", value)

		err = json.Unmarshal([]byte(valueString), &o)
		if err != nil {
			return nil, err
		}
	} else {
		valueByte, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(valueByte, &o)
		if err != nil {
			return nil, err
		}
	}

	if o.OrderUID == nil {
		return nil, fmt.Errorf("field [order_uid]  don't must be empty")
	}
	return
}
