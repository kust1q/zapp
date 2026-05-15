package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	cfg := &config{}
	err := cfg.Validate()
	assert.Error(t, err)

	cfg.App.Port = "8080"
	cfg.Elastic.Host = "localhost"
	cfg.Elastic.Port = "9200"
	cfg.GRPC.SearchPort = "9090"
	cfg.Kafka.Brokers = []string{"localhost:9092"}
	cfg.Kafka.Topics = []string{"topic"}
	cfg.Kafka.Consumer.GroupID = "group"

	err = cfg.Validate()
	assert.NoError(t, err)
}
