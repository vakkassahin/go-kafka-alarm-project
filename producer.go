package main

import (
	"fmt"
	"math/rand" // rand paketini ekleyin
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

// ... (diğer kodlar)

func main() {
	// Kafka broker'ına bağlan
	brokers := []string{"localhost:9092"} // Kafka broker adresi
	config := sarama.NewConfig()
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Printf("Error creating Kafka producer: %v\n", err)
		return
	}
	defer producer.Close()

	// Kafka topic'i belirle
	topic := "cpu_memory_metrics"

	// Ctrl+C sinyalini dinle
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Belirli aralıklarla sahte CPU ve bellek verileri üret
	go produceMetrics(producer, topic)

	// Kafka'dan mesajları tüket
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Printf("Error creating Kafka consumer: %v\n", err)
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		fmt.Printf("Error creating partition consumer: %v\n", err)
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			// Mesajları işle
			fmt.Printf("Received message: %s\n", msg.Value)
		case <-sigterm:
			// Uygulamayı kapat
			fmt.Println("Terminating...")
			return
		}
	}
}

func produceMetrics(producer sarama.SyncProducer, topic string) {
	for {
		cpuUsage := 80 + randFloat(-5, 5) // Sahte CPU kullanımı
		memoryUsage := 70 + randFloat(-3, 3) // Sahte bellek kullanımı

		message := fmt.Sprintf("CPU Usage: %.2f%%, Memory Usage: %.2f%%", cpuUsage, memoryUsage)
		// Kafka'ya mesajı gönder
		producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		})

		time.Sleep(5 * time.Second) // 5 saniyede bir metrik üret
	}
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
