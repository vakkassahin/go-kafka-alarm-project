package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

func main() {
	// Kafka broker'ına bağlan
	brokers := []string{"localhost:9092"} // Kafka broker adresi
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Printf("Error creating Kafka consumer: %v\n", err)
		return
	}
	defer consumer.Close()

	// Kafka topic'i belirle
	topic := "cpu_memory_metrics"

	// Ctrl+C sinyalini dinle
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	// Kafka'dan mesajları tüket ve işle
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
			handleMessage(msg.Value)
		case <-sigterm:
			// Uygulamayı kapat
			fmt.Println("Terminating...")
			return
		}
	}
}

func handleMessage(message []byte) {
	// Mesajı işle
	fmt.Printf("Received message: %s\n", message)

	// Burada metrikleri işleyerek belirli bir eşiği kontrol edebilir ve alarm üretebilirsiniz.
	// Örnek olarak, mesaj içinde CPU kullanımı kontrolü yapalım:
	if isCPUUsageHigh(message, 90.0) {
		fmt.Println("High CPU Usage Alarm!")
		// Burada bir alarm gönderme veya başka bir işlem yapabilirsiniz.
	}
}

func isCPUUsageHigh(message []byte, threshold float64) bool {
	// Mesajı analiz ederek CPU kullanımını kontrol et
	// Bu örnek basitçe mesajı string'e dönüştürüyor ve belirli bir eşiği kontrol ediyor.
	// Gerçek senaryoda mesaj formatınıza göre uyarlamanız gerekebilir.
	return parseCPUUsage(message) > threshold
}

func parseCPUUsage(message []byte) float64 {
	// Örnek olarak, mesaj içindeki CPU kullanımını alır.
	// Gerçek senaryoda mesaj formatınıza göre uyarlamanız gerekebilir.
	// Bu örnek mesajın tamamını alır ve başında "CPU Usage: " ifadesini arar.
	// Ardından, bu ifadenin sonrasındaki sayısal değeri (CPU kullanımını) pars eder.
	var cpuUsage float64
	fmt.Sscanf(string(message), "CPU Usage: %f%%", &cpuUsage)
	return cpuUsage
}
