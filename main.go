package main

import (
	"context"
	"time"

	"github.com/KScaesar/art"
	"github.com/segmentio/kafka-go"

	"github.com/KScaesar/GARMIN2024/pkg/infra"
)

var Topic_CreatedOrder = "Created.Order"
var Topic_V1_CreatedOrder = "v1." + Topic_CreatedOrder

func main() {
	art.SetDefaultLogger(art.NewLogger(false, art.LogLevelDebug))

	kafkaUrls := []string{"localhost:19092", "localhost:19093", "localhost:19094"}
	err := infra.PingKafka(kafkaUrls)
	if err != nil {
		panic(err)
	}

	producer, err := infra.KafkaFactory{
		ProducerName: "kafka_p",
		Writer:       infra.GeneralKafkaWriter(kafkaUrls),
		EgressMux:    NewKafkaEgressMux(),
	}.CreateProducer()
	if err != nil {
		panic(err)
	}
	go func() {
		egress := infra.NewBytesKafkaEgress(Topic_V1_CreatedOrder, []byte("hello_world"))
		infra.KafkaMetadata.SetOrderId(egress.Metadata, "xx1017-"+time.Now().String())
		err = producer.Send(egress)
		if err != nil {
			panic(err)
		}
	}()

	consumer, err := infra.KafkaFactory{
		ConsumerName: "kafka_c_v1.order",
		Reader:       infra.GeneralKafkaGroupReader(kafkaUrls, "v1.order", []string{Topic_V1_CreatedOrder}),
		IngressMux:   NewKafkaIngressMux(),
	}.CreateConsumer()
	if err != nil {
		panic(err)
	}
	go func() {
		consumer.Listen()
	}()

	art.NewShutdown().
		StopService("kafka_p", producer.Stop).
		StopService("kafka_c", consumer.Stop).
		Serve(context.Background())
}

func NewKafkaEgressMux() *infra.KafkaEgressMux {
	mux := infra.NewKafkaEgressMux().
		ErrorHandler(
			art.UsePrintResult{}.PrintEgress().PostMiddleware(),
		).
		Middleware(
			art.UseLogger(false, art.SafeConcurrency_Skip),
			infra.UseAutoCreateTopic(),
		).
		Handler(Topic_V1_CreatedOrder, infra.UseKafkaWriter(func(message *art.Message, w *kafka.Writer) error {
			orderId := infra.KafkaMetadata.GetOrderId(message.Metadata)
			return w.WriteMessages(message.Ctx, kafka.Message{
				Topic: message.Subject,
				Key:   []byte("OrderId=" + orderId),
				Value: message.Bytes,
			})
		}))
	return mux
}

func NewKafkaIngressMux() *infra.KafkaIngressMux {
	mux := infra.NewKafkaIngressMux().
		ErrorHandler(
			art.UsePrintResult{}.PrintIngress().PostMiddleware(),
		).
		Middleware(
			art.UseLogger(false, art.SafeConcurrency_Skip),
			art.UseAdHocFunc(func(message *art.Message, dep any) error {
				logger := art.CtxGetLogger(message.Ctx, dep)
				logger.Info("recv %v", message.Subject)
				return nil
			}).PreMiddleware(),
			infra.UseCommitMessageWhenHandleOk(),
		)

	mux.Group("v1.").
		Handler(Topic_CreatedOrder, art.UsePrintDetail())

	return mux
}
