package main

import (
	"log"
	"os"
    "time"
    "strconv"
    "strings"

	"github.com/Shopify/sarama"
)

func main() {
    brokerList := strings.Split(readEnv("KAFKA_BROKER_LIST", "localhost:9092"), ",")
    topic := readEnv("KAFKA_TOPIC", "input")
    maxRetry, _ := strconv.Atoi(readEnv("KAFKA_MAX_RETRY", strconv.Itoa(5)))

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

    for {
        now := time.Now()
        nanos := now.UnixNano()
        epochMs := nanos / 1000000

        msg := &sarama.ProducerMessage {
            Topic: topic,
            Value: sarama.StringEncoder(strconv.FormatInt(epochMs, 10)),
        }
        partition, offset, err := producer.SendMessage(msg)
        if err != nil {
            log.Panic(err)
        }
        log.Printf("Produced message: %s is stored in topic(%s)/partition(%d)/offset(%d)\n", msg.Value, topic, partition, offset)
        time.Sleep(time.Second * 1)
	}
}

func readEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}
