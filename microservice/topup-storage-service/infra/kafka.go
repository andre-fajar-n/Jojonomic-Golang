package infra

import (
	"log"
	"os"
	"strings"

	"github.com/segmentio/kafka-go"
)

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
