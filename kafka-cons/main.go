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
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true

	for i := 0; i < 5; i++ {
		go consumer(i, cfg)
	}

	time.Sleep(time.Minute)
}

func consumer(i int, cfg *sarama.Config) {

	csmr, err := sarama.NewConsumer([]string{kafkaAddr}, cfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer csmr.Close()

	partitionConsumer, err := csmr.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("%d Received message: %s\n", i, string(msg.Value))
		case err := <-partitionConsumer.Errors():
			log.Printf("%d Error: %v\n", i, err)
		}
	}
}
