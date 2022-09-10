package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertJsonToOrder_validValue(t *testing.T) {
	value := `{
	  "order_uid": "b563feb7b2b84b6test",
	  "track_number": "WBILMTESTTRACK",
	  "entry": "WBIL",
	  "locale": "en",
	  "internal_signature": "",
	  "customer_id": "test",
	  "delivery_service": "meest",
	  "shardkey": "9",
	  "sm_id": 99,
	  "date_created": "2021-11-26T06:22:19Z",
	  "oof_shard": "1"
	}`

	_, err := ConvertJsonToOrder(value)
	assert.Nil(t, err)
}

func TestConvertJsonToOrder_invalidValue(t *testing.T) {
	value := `{
	  "track_number": "WBILMTESTTRACK",
	  "entry": "WBIL",
	  "delivery": {},
	  "payment": {},
	  "items": [],
	  "locale": "en",
	  "internal_signature": "",
	  "customer_id": "test",
	  "delivery_service": "meest",
	  "shardkey": "9",
	  "sm_id": 99,
	  "date_created": "2021-11-26T06:22:19Z",
	  "oof_shard": "1"
	}`

	_, err := ConvertJsonToOrder(value)
	assert.NotNil(t, err)
}
