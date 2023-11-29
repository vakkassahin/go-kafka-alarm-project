package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Shopify/sarama"
)

// ConnectProducer is a function that creates a connection to Kafka
func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// PushMessage is a function that sends a message to Kafka
func PushMessage(conn sarama.SyncProducer, topic string, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := conn.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

// GetCPUUsage is a function that returns the current CPU usage percentage
func GetCPUUsage() (float64, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer file.Close()
	var user, nice, system, idle, iowait, irq, softirq, steal, guest, guestNice int
	_, err = fmt.Fscanf(file, "cpu %d %d %d %d %d %d %d %d %d %d",
		&user, &nice, &system, &idle, &iowait, &irq, &softirq, &steal, &guest, &guestNice)
	if err != nil {
		return 0, err
	}
	total := user + nice + system + idle + iowait + irq + softirq + steal
	usage := float64(total-idle) / float64(total) * 100
	return usage, nil
}

// GetMemoryUsage is a function that returns the current memory usage percentage
func GetMemoryUsage() (float64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()
	var total, free, available, buffers, cached int
	_, err = fmt.Fscanf(file, "MemTotal: %d kB\nMemFree: %d kB\nMemAvailable: %d kB\nBuffers: %d kB
