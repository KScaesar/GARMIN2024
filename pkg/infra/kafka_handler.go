package infra

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

//

func UseKafkaWriter(handler func(*art.Message, *kafka.Writer) error) art.HandleFunc {
	return func(message *art.Message, dependency any) error {
		dep := dependency.(KafkaProducer).RawInfra().(*kafka.Writer)
		return handler(message, dep)
	}
}

func UseAutoCreateTopic() art.Middleware {
	return func(next art.HandleFunc) art.HandleFunc {
		return func(message *art.Message, dependency any) error {
		retry:
			err := next(message, dependency)
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

func UseCommitMessageWhenHandleOk() art.Middleware {
	return func(next art.HandleFunc) art.HandleFunc {
		return func(message *art.Message, dependency any) error {
			err := next(message, dependency)
			if err != nil {
				return err
			}

			reader := dependency.(KafkaConsumer).RawInfra().(*kafka.Reader)
			return reader.CommitMessages(message.Ctx, *message.RawInfra.(*kafka.Message))
		}
	}
}
