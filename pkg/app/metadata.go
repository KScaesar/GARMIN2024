package app

import (
	"github.com/KScaesar/art"
)

var Metadata = newMetadataKey()

func newMetadataKey() *metadataKey {
	return &metadataKey{
		orderId: "orderId",
	}
}

type metadataKey struct {
	orderId string
}

func (key *metadataKey) GetOrderId(message *art.Message) string {
	md := message.Metadata
	return md.Get(key.orderId).(string)
}

func (key *metadataKey) SetOrderId(message *art.Message, value string) {
	md := message.Metadata
	md.Set(key.orderId, value)
}
