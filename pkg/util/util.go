package util

import (
	"encoding/json"
	"fmt"
	"github.com/Vitaly-Baidin/l0/pkg/entity"
	"reflect"
)

func ConvertJsonToOrder(value any) (o *entity.Order, err error) {
	valueType := reflect.TypeOf(value).Kind()
	if valueType == reflect.String {
		valueString := fmt.Sprintf("%v", value)

		err = json.Unmarshal([]byte(valueString), &o)
		if err != nil {
			return nil, fmt.Errorf("invalid message(string): %v", err)
		}

	} else if valueType == reflect.Slice {
		err = json.Unmarshal(value.([]byte), &o)
	} else {
		fmt.Println(valueType)
		valueByte, err := json.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("invalid message(order): %v", err)
		}

		err = json.Unmarshal(valueByte, &o)
		if err != nil {
			return nil, fmt.Errorf("invalid message(order): %v", err)
		}
	}

	if o.OrderUID == nil {
		return nil, fmt.Errorf("field [order_uid] don't must be empty")
	}

	return
}
