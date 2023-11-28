package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Kafka struct {
		Broker string `mapstructure:"broker"`
		Topic  string `mapstructure:"topic"`
	} `mapstructure:"kafka"`

	Logging struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logging"`
}

func main() {
	// Viper konfigürasyon yükleyiciyi oluştur
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Konfigürasyon dosyasını oku
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}

	// Konfigürasyon yapısını tanımla
	var config Config

	// Konfigürasyonu yükle
	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Error unmarshaling config: %v\n", err)
		return
	}

	// Konfigürasyonu kullan
	fmt.Printf("Kafka Broker: %s\n", config.Kafka.Broker)
	fmt.Printf("Kafka Topic: %s\n", config.Kafka.Topic)
	fmt.Printf("Logging Level: %s\n", config.Logging.Level)
}
