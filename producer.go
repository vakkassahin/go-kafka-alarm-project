package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
)

const (
	kafkaBrokers = "localhost:9092"
	topic        = "temperature_data"
	alarmTopic   = "temperature_alarms"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewSyncProducer([]string{kafkaBrokers}, config)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				temperature := rand.Float64() * 100.0
				fmt.Printf("Current temperature: %.2f\n", temperature)

				// Send temperature data to Kafka
				msg := &sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(fmt.Sprintf("%.2f", temperature)),
				}
				_, _, err := producer.SendMessage(msg)
				if err != nil {
					log.Printf("Failed to produce message: %s", err)
				}

				// Check for alarm condition
				if temperature > 90.0 {
					// Send alarm message to Kafka
					alarmMsg := &sarama.ProducerMessage{
						Topic: alarmTopic,
						Value: sarama.StringEncoder(fmt.Sprintf("High temperature alarm! %.2f", temperature)),
					}
					_, _, err := producer.SendMessage(alarmMsg)
					if err != nil {
						log.Printf("Failed to produce alarm message: %s", err)
					}
				}
			case <-doneCh:
				return
			}
		}
	}()

	<-signals
	close(doneCh)
	fmt.Println("Exiting...")
}
