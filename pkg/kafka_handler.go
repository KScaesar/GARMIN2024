package pkg

import (
	"errors"
	"time"

	"github.com/KScaesar/art"
	"github.com/segmentio/kafka-go"
)

type KafkaIngressMux = art.Mux
type KafkaEgressMux = art.Mux

func NewKafkaIngressMux() *KafkaIngressMux {
	in := art.NewMux(".")
	return in
}

func NewKafkaEgressMux() *KafkaEgressMux {
	out := art.NewMux(".")
	return out
}

// consumer handler

func kafkaIngress(kafkaMessage kafka.Message) *art.Message {
	message := art.GetMessage()

	message.Subject = kafkaMessage.Topic
	message.Bytes = kafkaMessage.Value
	message.RawInfra = kafkaMessage
	return message
}

func UseCommitKafkaWhenHandleOk() art.Middleware {
	return func(next art.HandleFunc) art.HandleFunc {
		return func(ingress *art.Message, dependency any) error {
			err := next(ingress, dependency)
			if err != nil {
				return err
			}

			reader := dependency.(KafkaConsumer).RawInfra().(*kafka.Reader)
			return reader.CommitMessages(ingress.Ctx, ingress.RawInfra.(kafka.Message))
		}
	}
}

// producer handler

func UseKafkaWriter(handler func(*art.Message, *kafka.Writer) error) art.HandleFunc {
	return func(egress *art.Message, dependency any) error {
		dep := dependency.(KafkaProducer).RawInfra().(*kafka.Writer)
		return handler(egress, dep)
	}
}

func UseAutoCreateKafkaTopic() art.Middleware {
	return func(next art.HandleFunc) art.HandleFunc {
		return func(egress *art.Message, dependency any) error {
		retry:
			err := next(egress, dependency)
			if err != nil {
				if errors.Is(err, kafka.UnknownTopicOrPartition) {
					time.Sleep(500 * time.Millisecond)
					goto retry
				}
			}
			return err
		}
	}
}
