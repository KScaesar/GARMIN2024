package main

import (
	"context"
	"net/http"

	"github.com/KScaesar/art"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/KScaesar/GARMIN2024/pkg"
	"github.com/KScaesar/GARMIN2024/pkg/adapters/api"
	"github.com/KScaesar/GARMIN2024/pkg/adapters/pubsub"
	"github.com/KScaesar/GARMIN2024/pkg/app"
)

func main() {
	art.SetDefaultLogger(art.NewLogger(false, art.LogLevelDebug))
	logger := art.DefaultLogger()

	conf, err := pkg.LoadJsonConfingByLocal("")
	if err != nil {
		logger.Fatal(err.Error())
	}

	err = pkg.PingKafka(conf.KafkaUrls)
	if err != nil {
		logger.Fatal(err.Error())
	}

	producer, err := pubsub.NormalKafkaProducer(&conf)
	if err != nil {
		logger.Fatal(err.Error())
	}

	svc := &app.Service{
		OrderService:    app.NewOrderUseCase(producer),
		ShippingService: app.NewShippingUseCase(),
	}

	router := api.NewGinRouter(&conf, svc)
	httpServer := api.NewHttpServer(&conf, router)
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	consumer, err := pubsub.NormalKafkaConsumer(&conf, svc)
	if err != nil {
		logger.Fatal(err.Error())
	}
	go func() {
		consumer.Listen()
	}()

	ServeO11Y(&conf.O11Y)

	logger.Info("service started")
	ctx := context.Background()
	art.NewShutdown().
		StopService("http", func() error {
			return httpServer.Shutdown(ctx)
		}).
		StopService("kafka_p", producer.Stop).
		StopService("kafka_c", consumer.Stop).
		Serve(ctx)
}

func ServeO11Y(conf *pkg.O11Y) {
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":"+conf.MetricPort, http.DefaultServeMux)
	}()
}
