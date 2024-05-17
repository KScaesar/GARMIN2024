package infra

import (
	"github.com/gookit/goutil/maputil"
)

var KafkaMetadata = newKafkaMetadataKey()

func newKafkaMetadataKey() *kafkaMetadataKey {
	return &kafkaMetadataKey{
		orderId: "orderId",
	}
}

type kafkaMetadataKey struct {
	orderId string
}

func (key *kafkaMetadataKey) GetOrderId(md maputil.Data) string {
	return md.Get(key.orderId).(string)
}

func (key *kafkaMetadataKey) SetOrderId(md maputil.Data, value string) {
	md.Set(key.orderId, value)
}
