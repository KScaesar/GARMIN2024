package main

import (
	"context"
	"errors"

	"github.com/segmentio/kafka-go"
)

func main() {
	kafkaUrls := []string{"localhost:19092", "localhost:19093", "localhost:19094"}
	err := PingKafka(kafkaUrls)
	if err != nil {
		panic(err)
	}

	writer := GeneralKafkaWriter(kafkaUrls)
retry:
	err = writer.WriteMessages(context.Background(), kafka.Message{
		Topic: "Created.Order",
		Key:   []byte("OrderId=" + "f1017"),
		Value: []byte("Hello World"),
	})
	if err != nil {
		if errors.Is(err, kafka.UnknownTopicOrPartition) {
			goto retry
		}
		panic(err)
	}
}

func PingKafka(urls []string) error {
	for _, url := range urls {
		conn, err := kafka.DefaultDialer.Dial("tcp", url)
		if err != nil {
			return err
		}

		err = conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func GeneralKafkaWriter(urls []string) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(urls...),
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true,
	}
	return writer
}
