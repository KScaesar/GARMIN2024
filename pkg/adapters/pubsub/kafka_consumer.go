package pubsub

import (
	"github.com/KScaesar/art"

	"github.com/KScaesar/GARMIN2024/pkg"
	"github.com/KScaesar/GARMIN2024/pkg/app"
)

func NormalKafkaConsumer(conf *pkg.Config, svc *app.Service) (pkg.KafkaConsumer, error) {
	err := pkg.CreateKafkaTopic(conf.KafkaUrls[0], []string{
		app.Subject_V1_CreatedOrder,
	})
	if err != nil {
		return nil, err
	}

	return pkg.KafkaFactory{
		ConsumerName: "kafka_c_v1.order",
		Reader: pkg.NormalKafkaGroupReader(
			conf.KafkaUrls,
			"v1.order",
			[]string{app.Subject_V1_CreatedOrder},
		),
		IngressMux: kafkaIngressMux(svc),
	}.CreateConsumer()
}

func kafkaIngressMux(svc *app.Service) *pkg.KafkaIngressMux {
	mux := pkg.NewKafkaIngressMux().
		ErrorHandler(
			art.UsePrintResult{}.PrintIngress().PostMiddleware(),
		).
		Middleware(
			art.UseLogger(false, art.SafeConcurrency_Skip),
			art.UseAdHocFunc(func(message *art.Message, dep any) error {
				logger := art.CtxGetLogger(message.Ctx)
				logger.Info("kafka recv %q", message.Subject)
				return nil
			}).PreMiddleware(),
			pkg.UseCommitKafkaWhenHandleOk(),
		)

	mux.Group("v1.").
		Handler(app.Subject_CreatedOrder, createShipping(svc))

	return mux
}
