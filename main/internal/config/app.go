package config

import (
	"crypto/rsa"
	"time"
)

type (
	ApplicationConfig struct {
		Port string `mapstructure:"port"`
	}

	CacheConfig struct {
		DefaultTtl  time.Duration `mapstructure:"default_ttl"`
		CountersTtl time.Duration `mapstructure:"counters_ttl"`
	}

	TokensConfig struct {
		AccessTTL   time.Duration `mapstructure:"access_ttl"`
		RefreshTTL  time.Duration `mapstructure:"refresh_ttl"`
		RecoveryTTL time.Duration `mapstructure:"recovery_ttl"`
	}

	JWTConfig struct {
		PrivateKey *rsa.PrivateKey
		PublicKey  *rsa.PublicKey
	}

	MinioConfig struct {
		Port       string        `mapstructure:"port"`
		Endpoint   string        `mapstructure:"endpoint"`
		BucketName string        `mapstructure:"bucketname"`
		User       string        `mapstructure:"user"`
		Password   string        `mapstructure:"password"`
		UseSSL     bool          `mapstructure:"sslmode"`
		TTL        time.Duration `mapstructure:"ttl"`
	}

	PostgresConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"db_name"`
		SSLMode  string `mapstructure:"sslmode"`
	}

	RedisConfig struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}

	ElasticConfig struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	GrpcConfig struct {
		Host            string `mapstructure:"host"`
		IntegrationPort string `mapstructure:"integration_port"`
		SearchPort      string `mapstructure:"search_port"`
	}

	KafkaConfig struct {
		Brokers  []string `mapstructure:"brokers"`
		Topics   []string `mapstructure:"topics"`
		Producer struct {
			MaxRetries int `mapstructure:"max_retries"`
		} `mapstructure:"producer"`
		Consumer struct {
			GroupID string `mapstructure:"group_id"`
		} `mapstructure:"consumer"`
	}
)
