package kafka

import (
	"fmt"
	"os"

	confluentKafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewKafkaProducer() *confluentKafka.Producer {
	configmap := &confluentKafka.ConfigMap{
		"bootstrap.servers": os.Getenv("kafkaBootstrapServers"),
	}
	p, err := confluentKafka.NewProducer(configmap)
	if err != nil {
		panic(err)
	}
	return p
}

func Publish(msg string, topic string, producer *confluentKafka.Producer, deliveryChan chan confluentKafka.Event) error {
	message := confluentKafka.Message{
		TopicPartition: confluentKafka.TopicPartition{Topic: &topic, Partition: confluentKafka.PartitionAny},
		Value:          []byte(msg),
	}
	return producer.Produce(&message, deliveryChan)
}

func DeliveryReport(deliveryChan chan confluentKafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *confluentKafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Delivery failed:", ev.TopicPartition)
			} else {
				fmt.Println("Delivered to:", ev.TopicPartition)
			}
		}
	}
}
