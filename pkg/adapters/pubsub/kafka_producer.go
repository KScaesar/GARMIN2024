package pubsub

import (
	"time"

	"github.com/KScaesar/art"
	"github.com/segmentio/kafka-go"

	"github.com/KScaesar/GARMIN2024/pkg"
	"github.com/KScaesar/GARMIN2024/pkg/app"
)

func NormalKafkaProducer(conf *pkg.Config) (pkg.KafkaProducer, error) {
	var w *kafka.Writer
	if conf.Production {
		w = pkg.NormalKafkaWriter(conf.KafkaUrls, time.Second)
	} else {
		w = pkg.NormalKafkaWriter(conf.KafkaUrls, time.Millisecond)
	}

	return pkg.KafkaFactory{
		ProducerName: "kafka_p",
		Writer:       w,
		EgressMux:    kafkaEgressMux(),
	}.CreateProducer()
}

func kafkaEgressMux() *pkg.KafkaEgressMux {
	mux := pkg.NewKafkaEgressMux().
		ErrorHandler(
			art.UsePrintResult{}.PrintEgress().PostMiddleware(),
		).
		Middleware(
			art.UseLogger(true, art.SafeConcurrency_Skip),
			art.UsePrintDetail().PreMiddleware(),
			pkg.UseEncodeJson(),
			pkg.UseAutoCreateKafkaTopic(),
		).
		Handler(app.Subject_V1_CreatedOrder, pkg.UseKafkaWriter(createdOrder))
	return mux
}
