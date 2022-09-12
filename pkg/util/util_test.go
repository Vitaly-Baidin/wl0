package util

import (
	"encoding/json"
	"github.com/Vitaly-Baidin/l0/pkg/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConvertJsonToOrder_stringValue(t *testing.T) {
	t.Parallel()
	t.Run("test String value", func(t *testing.T) {
		t.Parallel()
		t.Run("valid value", func(t *testing.T) {
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
			require.NoError(t, err)
		})
		t.Run("invalid value", func(t *testing.T) {
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
			require.Error(t, err)
		})
	})
}

func TestConvertJsonToOrder_orderStructValue(t *testing.T) {
	t.Parallel()
	t.Run("test Order value", func(t *testing.T) {
		t.Parallel()
		t.Run("valid value", func(t *testing.T) {
			orderUID := "123abc"
			value := entity.Order{
				OrderUID:          &orderUID,
				TrackNumber:       "",
				Entry:             "",
				DeliveryData:      nil,
				PaymentData:       nil,
				ItemsData:         nil,
				Locale:            "",
				InternalSignature: "",
				CustomerID:        "",
				DeliveryService:   "",
				Shardkey:          "",
				SmID:              0,
				DateCreated:       time.Time{},
				OofShard:          "",
			}

			_, err := ConvertJsonToOrder(value)
			require.NoError(t, err)
		})
		t.Run("invalid value", func(t *testing.T) {
			var orderUID string
			value := entity.Order{
				OrderUID:          &orderUID,
				TrackNumber:       "",
				Entry:             "",
				DeliveryData:      nil,
				PaymentData:       nil,
				ItemsData:         nil,
				Locale:            "",
				InternalSignature: "",
				CustomerID:        "",
				DeliveryService:   "",
				Shardkey:          "",
				SmID:              0,
				DateCreated:       time.Time{},
				OofShard:          "",
			}

			_, err := ConvertJsonToOrder(value)
			require.Error(t, err)
		})
	})
}

func TestConvertJsonToOrder_sliceByteValue(t *testing.T) {
	t.Parallel()
	t.Run("test Order value", func(t *testing.T) {
		t.Parallel()
		t.Run("valid value", func(t *testing.T) {
			orderUID := "123abc"
			value := entity.Order{
				OrderUID:          &orderUID,
				TrackNumber:       "",
				Entry:             "",
				DeliveryData:      nil,
				PaymentData:       nil,
				ItemsData:         nil,
				Locale:            "",
				InternalSignature: "",
				CustomerID:        "",
				DeliveryService:   "",
				Shardkey:          "",
				SmID:              0,
				DateCreated:       time.Time{},
				OofShard:          "",
			}

			byteValue, err := json.Marshal(value)
			assert.NoError(t, err)

			_, err = ConvertJsonToOrder(byteValue)
			require.NoError(t, err)
		})
		t.Run("invalid value", func(t *testing.T) {
			var orderUID string
			value := entity.Order{
				OrderUID:          &orderUID,
				TrackNumber:       "",
				Entry:             "",
				DeliveryData:      nil,
				PaymentData:       nil,
				ItemsData:         nil,
				Locale:            "",
				InternalSignature: "",
				CustomerID:        "",
				DeliveryService:   "",
				Shardkey:          "",
				SmID:              0,
				DateCreated:       time.Time{},
				OofShard:          "",
			}

			byteValue, err := json.Marshal(value)
			assert.NoError(t, err)

			_, err = ConvertJsonToOrder(byteValue)
			require.Error(t, err)
		})
	})
}

func TestConvertJsonToOrder_invalidValue(t *testing.T) {
	t.Parallel()
	t.Run("test any string value", func(t *testing.T) {
		_, err := ConvertJsonToOrder("sdasd")
		require.Error(t, err)
	})
	t.Run("test any int value", func(t *testing.T) {
		_, err := ConvertJsonToOrder(1234541)
		require.Error(t, err)
	})
	t.Run("test any bool value", func(t *testing.T) {
		_, err := ConvertJsonToOrder(true)
		require.Error(t, err)
	})
	t.Run("test any slice byte value", func(t *testing.T) {
		byteValue, err := json.Marshal("asdasd")
		assert.NoError(t, err)
		_, err = ConvertJsonToOrder(byteValue)
		require.Error(t, err)
	})
}
