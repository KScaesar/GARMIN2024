package infra

import (
	"github.com/KScaesar/art"
	"github.com/segmentio/kafka-go"
)

func NewKafkaIngress(kafkaMessage *kafka.Message) *art.Message {
	message := art.GetMessage()

	message.Subject = kafkaMessage.Topic
	message.Bytes = kafkaMessage.Value
	message.RawInfra = kafkaMessage
	return message
}

func NewBytesKafkaEgress(topic string, bMessage []byte) *art.Message {
	message := art.GetMessage()

	message.Subject = topic
	message.Bytes = bMessage
	return message
}

func NewBodyKafkaEgress(topic string, body any) *art.Message {
	message := art.GetMessage()

	message.Subject = topic
	message.Body = body
	return message
}
