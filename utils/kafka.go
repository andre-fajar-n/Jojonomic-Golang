package utils

import (
	"log"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

func GetKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	l := log.New(os.Stdout, "kafka writer: ", 0)
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Logger:   l,
	}
}

func GetKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	l := log.New(os.Stdout, "kafka reader: ", 0)
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		Logger:   l,
	})
}
