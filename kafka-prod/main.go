package main

import (
	"github.com/IBM/sarama"
	"log"
	"time"
)

const (
	kafkaAddr = "10.161.33.38:9092"
	topic     = "test-topic"
)

func main() {
	producer()
}

func producer() {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	prod, err := sarama.NewSyncProducer([]string{kafkaAddr}, cfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer prod.Close()

	for i := 0; i < 10; i++ {
		msg := sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder("Hello Kafka!")}

		partition, offset, err := prod.SendMessage(&msg)
		if err != nil {
			log.Fatalf("err: %v", err)
		}
		log.Printf("Message sent to partition %d at offset %d\n", partition, offset)

		time.Sleep(time.Second)
	}
}
