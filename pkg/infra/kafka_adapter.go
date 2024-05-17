package infra

import (
	"context"

	"github.com/KScaesar/art"
	"github.com/segmentio/kafka-go"
)

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

func GeneralKafkaGroupReader(urls []string, groupId string, topics []string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     urls,
		GroupID:     groupId,
		GroupTopics: topics,
	})
}

//

type KafkaProducer = art.Producer
type KafkaConsumer = art.Consumer

//

type KafkaFactory struct {
	Hub    *art.Hub
	Logger art.Logger

	ProducerName string
	Writer       *kafka.Writer
	EgressMux    *KafkaEgressMux

	ConsumerName string
	Reader       *kafka.Reader
	IngressMux   *KafkaIngressMux

	DecorateAdapter func(adapter art.IAdapter) (application art.IAdapter)
	Lifecycle       func(lifecycle *art.Lifecycle)
}

func (f KafkaFactory) CreateProducer() (producer KafkaProducer, err error) {
	opt := art.NewAdapterOption().
		Identifier(f.ProducerName).
		RawInfra(f.Writer).
		EgressMux(f.EgressMux).
		AdapterHub(f.Hub).
		Logger(f.Logger).
		DecorateAdapter(f.DecorateAdapter).
		Lifecycle(f.Lifecycle)

	opt.RawStop(func(logger art.Logger) (err error) {
		err = f.Writer.Close()
		if err != nil {
			logger.Error("writer stop: %v", err)
			return
		}
		return nil
	})

	adp, err := opt.Build()
	if err != nil {
		return
	}
	return adp.(KafkaProducer), err
}

func (f KafkaFactory) CreateConsumer() (consumer KafkaConsumer, err error) {
	opt := art.NewAdapterOption().
		Identifier(f.ConsumerName).
		RawInfra(f.Reader).
		IngressMux(f.IngressMux).
		AdapterHub(f.Hub).
		Logger(f.Logger).
		DecorateAdapter(f.DecorateAdapter).
		Lifecycle(f.Lifecycle)

	opt.RawRecv(func(logger art.Logger) (message *art.Message, err error) {
		kafkaMessage, err := f.Reader.FetchMessage(context.Background())
		if err != nil {
			logger.Error("recv: %v", err)
			return nil, err
		}
		return NewKafkaIngress(&kafkaMessage), nil
	})

	opt.RawStop(func(logger art.Logger) (err error) {
		err = f.Reader.Close()
		if err != nil {
			logger.Error("reader stop: %v", err)
			return
		}
		return nil
	})

	adp, err := opt.Build()
	if err != nil {
		return
	}
	return adp.(KafkaConsumer), err
}
