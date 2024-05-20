package pubsub

import (
	"github.com/KScaesar/art"
	"github.com/segmentio/kafka-go"

	"github.com/KScaesar/GARMIN2024/pkg/app"
)

func createdOrder(egress *art.Message, w *kafka.Writer) error {
	orderId := app.Metadata.GetOrderId(egress)

	return w.WriteMessages(egress.Ctx, kafka.Message{
		Topic: app.Subject_V1_CreatedOrder,
		Key:   []byte("OrderId=" + orderId),
		Value: egress.Bytes,
	})
}
