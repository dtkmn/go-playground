package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func init() {
	// Initialize Kafka writer
	kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKERS")},
		Topic:    "my-topic",
		Balancer: &kafka.LeastBytes{},
	})
}

func sendMessage(message string) {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Send message to Kafka
	err := kafkaWriter.WriteMessages(ctx, kafka.Message{Value: []byte(message)})
	if err != nil {
		log.Fatal(err)
	}
}
