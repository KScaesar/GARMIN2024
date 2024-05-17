package pkg

import (
	"context"
	"time"

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

func CreateKafkaTopic(url string, topics []string) error {
	conn, err := kafka.Dial("tcp", url)
	if err != nil {
		return err
	}
	defer conn.Close()

	conf := []kafka.TopicConfig{}
	for _, topic := range topics {
		conf = append(conf, kafka.TopicConfig{
			Topic:              topic,
			NumPartitions:      3,
			ReplicationFactor:  2,
			ReplicaAssignments: nil,
			ConfigEntries:      nil,
		})
	}
	return conn.CreateTopics(conf...)
}

func NormalKafkaWriter(urls []string, batchTimeout time.Duration) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(urls...),
		Balancer:               &kafka.Hash{},
		BatchTimeout:           batchTimeout,
		AllowAutoTopicCreation: true,
	}
	return writer
}

func NormalKafkaGroupReader(urls []string, groupId string, topics []string) *kafka.Reader {
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
			logger.Error("kafka writer stop: %v", err)
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
			logger.Error("kafka recv: %v", err)
			return nil, err
		}
		return kafkaIngress(kafkaMessage), nil
	})

	opt.RawStop(func(logger art.Logger) (err error) {
		err = f.Reader.Close()
		if err != nil {
			logger.Error("kafka reader stop: %v", err)
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
