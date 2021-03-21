package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var (
    brokerList = strings.Split(readEnv("KAFKA_BROKER_LIST", "localhost:9092"), ",")
)

func main() {
    consumerTopic := readEnv("KAFKA_CONSUMER_TOPIC", "input")
    consumerMessageCountStart, _ := strconv.Atoi(readEnv("KAFKA_CONSUMER_MESSAGE_COUNT_START", strconv.Itoa(1)))

	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_0
	config.Consumer.Return.Errors = true
	brokers := brokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(consumerTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				consumerMessageCountStart++
				log.Printf("Received message: %s", string(msg.Value))
				epochMs, _ := strconv.ParseInt(string(msg.Value), 10, 64)
				produce(epochMs)
			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed", consumerMessageCountStart, "messages")
}


func produce(epochMs int64) {
    producerTopic := readEnv("KAFKA_PRODUCER_TOPIC", "output")
    maxRetry, _ := strconv.Atoi(readEnv("KAFKA_PRODUCER_MAX_RETRY", strconv.Itoa(5)))

	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = maxRetry
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()

    epochS := epochMs / 1000
    msg := &sarama.ProducerMessage{
      Topic: producerTopic,
      Value: sarama.StringEncoder(time.Unix(epochS, 0).Format(time.RFC3339)),
    }
    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
      log.Panic(err)
    }
    log.Printf("Produced message: %s is stored in topic(%s)/partition(%d)/offset(%d)\n", msg.Value, producerTopic, partition, offset)
}

func readEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}