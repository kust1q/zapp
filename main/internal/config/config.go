package config

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type config struct {
	App      ApplicationConfig `mapstructure:"app"`
	Postgres PostgresConfig    `mapstructure:"db"`
	Minio    MinioConfig       `mapstructure:"minio"`
	Redis    RedisConfig       `mapstructure:"redis"`
	Elastic  ElasticConfig     `mapstructure:"elastic"`
	Cache    CacheConfig       `mapstructure:"cache"`
	Tokens   TokensConfig      `mapstructure:"tokens"`
	GRPC     GrpcConfig        `mapstructure:"grpc"`
	Kafka    KafkaConfig       `mapstructure:"kafka"`
	JWT      JWTConfig
}

var (
	instance *config
	once     sync.Once
)

func Get() *config {
	privateKey, publicKey, err := loadRSAKeys(os.Getenv("PRIVATE_KEY_PATH"), os.Getenv("PUBLIC_KEY_PATH"))
	if err != nil {
		logrus.Fatalf("Failed to load RSA keys: %v", err)
	}

	once.Do(func() {
		var cfg config
		if err := viper.Unmarshal(&cfg); err != nil {
			logrus.Fatalf("viper unmarshal failed: %v", err)
		}

		cfg.Postgres.User = os.Getenv("POSTGRES_USER")
		cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
		cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")
		cfg.Minio.User = os.Getenv("MINIO_USER")
		cfg.Minio.Password = os.Getenv("MINIO_PASSWORD")
		cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")

		cfg.JWT.PrivateKey = privateKey
		cfg.JWT.PublicKey = publicKey

		instance = &cfg
	})
	return instance
}

func (c *config) Validate() error {
	var allErrs []string

	if c.App.Port == "" {
		allErrs = append(allErrs, "app: port is required")
	}

	pg := c.Postgres
	validModes := map[string]bool{
		"disable": true, "allow": true, "prefer": true,
		"require": true, "verify-ca": true, "verify-full": true,
	}
	if pg.Host == "" {
		allErrs = append(allErrs, "postgres: host is required")
	}
	if pg.Port == "" {
		allErrs = append(allErrs, "postgres: port is required")
	}
	if pg.User == "" {
		allErrs = append(allErrs, "postgres: user is required")
	}
	if pg.Password == "" {
		allErrs = append(allErrs, "postgres: password is required")
	}
	if pg.DBName == "" {
		allErrs = append(allErrs, "postgres: dbname is required")
	}
	if !validModes[pg.SSLMode] {
		allErrs = append(allErrs, "postgres: invalid sslmode")
	}

	minio := c.Minio
	if minio.Port == "" {
		allErrs = append(allErrs, "minio: port is required")
	}
	if minio.Endpoint == "" {
		allErrs = append(allErrs, "minio: endpoint is required")
	}
	if minio.BucketName == "" {
		allErrs = append(allErrs, "minio: bucketname is required")
	}
	if minio.User == "" {
		allErrs = append(allErrs, "minio: user is required")
	}
	if minio.Password == "" {
		allErrs = append(allErrs, "minio: password is required")
	}
	if minio.TTL <= 0 {
		allErrs = append(allErrs, "minio: ttl is required")
	}

	redis := c.Redis
	if redis.Host == "" {
		allErrs = append(allErrs, "redis: host is required")
	}
	if redis.Port == "" {
		allErrs = append(allErrs, "redis: port is required")
	}
	if redis.DB < 0 {
		allErrs = append(allErrs, "redis: db must be >= 0")
	}

	el := c.Elastic
	if el.Host == "" {
		allErrs = append(allErrs, "elastic: host is required")
	}
	if el.Port == "" {
		allErrs = append(allErrs, "elastic: port is required")
	}

	if c.Cache.DefaultTtl <= 0 {
		allErrs = append(allErrs, "cache: default ttl must be > 0")
	}

	if c.Cache.CountersTtl <= 0 {
		allErrs = append(allErrs, "cache: counters ttl must be > 0")
	}

	if c.Tokens.AccessTTL <= 0 {
		allErrs = append(allErrs, "tokens: access ttl must be > 0")
	}
	if c.Tokens.RefreshTTL <= 0 {
		allErrs = append(allErrs, "tokens: refresh ttl must be > 0")
	}
	if c.Tokens.RecoveryTTL <= 0 {
		allErrs = append(allErrs, "tokens: recovery ttl must be > 0")
	}
	if c.JWT.PrivateKey == nil {
		allErrs = append(allErrs, "jwt: private key path is required")
	}
	if c.JWT.PublicKey == nil {
		allErrs = append(allErrs, "jwt: public key path is required")
	}

	if c.GRPC.Host == "" {
		allErrs = append(allErrs, "grpc: host is required")
	}
	if c.GRPC.SearchPort == "" {
		allErrs = append(allErrs, "grpc: search port is required")
	}
	if c.GRPC.IntegrationPort == "" {
		allErrs = append(allErrs, "grpc: integration port is required")
	}

	if len(c.Kafka.Brokers) == 0 {
		allErrs = append(allErrs, "kafka: brokers is required")
	}
	if len(c.Kafka.Topics) == 0 {
		allErrs = append(allErrs, "kafka: events names is required")
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
		logrus.Warnf(".env file not loaded: %v", err)
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

func loadRSAKeys(privateKeyPath, publicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privatePEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read private key from %s: %w", privateKeyPath, err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	return privateKey, publicKey, nil
}
