package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/doublegrey/microlog/utils"
)

func Produce(message []byte, topic string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": utils.Config.Kafka.Brokers})
	if err != nil {
		return err
	}

	delivery := make(chan kafka.Event)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}, delivery)
	// TODO: implement timeout
	event := <-delivery
	close(delivery)
	response := event.(*kafka.Message)
	if response.TopicPartition.Error != nil {
		return response.TopicPartition.Error
	}
	return nil
}

// FIXME: remove infinite loop (wait for sigterm or some event from main function to stop)
func Consume(topics []string, messages chan *kafka.Message, offset string) error {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     utils.Config.Kafka.Brokers,
		"broker.address.family": "v4",
		"group.id":              "0",
		"session.timeout.ms":    6000,
		"auto.offset.reset":     offset})
	if err != nil {
		return err
	}
	err = c.SubscribeTopics(topics, nil)
	for {
		ev := c.Poll(100)
		if ev == nil {
			continue
		}
		switch e := ev.(type) {
		case *kafka.Message:
			fmt.Printf("%% Message on %s:\n%s\n",
				e.TopicPartition, string(e.Value))
			if e.Headers != nil {
				fmt.Printf("%% Headers: %v\n", e.Headers)
			}
			messages <- e
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			if e.Code() == kafka.ErrAllBrokersDown {
				return errors.New("All brokers are down")
			}
		}
	}
}
