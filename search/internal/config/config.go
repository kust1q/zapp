package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type config struct {
	App     ApplicationConfig `mapstructure:"app"`
	Elastic ElasticConfig     `mapstructure:"elastic"`
	GRPC    GrpcConfig        `mapstructure:"grpc"`
	Kafka   KafkaConfig       `mapstructure:"kafka"`
}

var (
	instance *config
	once     sync.Once
)

func Get() *config {
	once.Do(func() {
		var cfg config
		if err := viper.Unmarshal(&cfg); err != nil {
			logrus.Fatalf("viper unmarshal failed: %v", err)
		}
		instance = &cfg
	})
	return instance
}

func (c *config) Validate() error {
	var allErrs []string

	if c.App.Port == "" {
		allErrs = append(allErrs, "app: port is required")
	}

	el := c.Elastic
	if el.Host == "" {
		allErrs = append(allErrs, "elastic: host is required")
	}
	if el.Port == "" {
		allErrs = append(allErrs, "elastic: port is required")
	}

	if c.GRPC.SearchPort == "" {
		allErrs = append(allErrs, "grpc: search port is required")
	}

	if len(c.Kafka.Brokers) == 0 {
		allErrs = append(allErrs, "kafka: brokers is required")
	}
	if len(c.Kafka.Topics) == 0 {
		allErrs = append(allErrs, "kafka: topics is required")
	}
	if c.Kafka.Consumer.GroupID == "" {
		allErrs = append(allErrs, "kafka: consumer group id is required")
	}

	if len(allErrs) > 0 {
		return errors.New("config validation errors: " + strings.Join(allErrs, " "))
	}
	return nil
}

func InitConfig() error {
	viper.SetConfigFile(".env")
	if err := viper.MergeInConfig(); err != nil {
		logrus.Warnf(".env file loaded: %v", err)
	}

	viper.AutomaticEnv()
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, assuming env vars are set")
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("configs")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
}
