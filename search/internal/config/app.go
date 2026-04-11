package config

type (
	ApplicationConfig struct {
		Port string `mapstructure:"port"`
	}

	ElasticConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	GrpcConfig struct {
		SearchPort string `mapstructure:"search_port"`
	}

	KafkaConfig struct {
		Brokers  []string `mapstructure:"brokers"`
		Topics   []string `mapstructure:"topics"`
		Consumer struct {
			GroupID string `mapstructure:"group_id"`
		} `mapstructure:"consumer"`
	}
)
