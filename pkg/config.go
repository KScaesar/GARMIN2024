package pkg

type Config struct {
	Production bool
	HttpPort   string
	GinDebug   bool
	KafkaUrls  []string
}
