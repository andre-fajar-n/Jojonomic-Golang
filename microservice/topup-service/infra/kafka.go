package infra

import (
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func GetKafkaWriter(topic string) *kafka.Writer {
	l := log.New(os.Stdout, "kafka writer: ", 0)
	return &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("kafka_url")),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Logger:   l,
	}
}
