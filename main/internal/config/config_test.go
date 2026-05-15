package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate(t *testing.T) {
	cfg := &config{}
	err := cfg.Validate()
	assert.Error(t, err)

	cfg.App.Port = "8080"
	cfg.Postgres.Host = "localhost"
	cfg.Postgres.Port = "5432"
	cfg.Postgres.User = "user"
	cfg.Postgres.Password = "pass"
	cfg.Postgres.DBName = "db"
	cfg.Postgres.SSLMode = "disable"
	
	cfg.Minio.Port = "9000"
	cfg.Minio.Endpoint = "localhost"
	cfg.Minio.BucketName = "bucket"
	cfg.Minio.User = "user"
	cfg.Minio.Password = "pass"
	cfg.Minio.TTL = time.Hour

	cfg.Redis.Host = "localhost"
	cfg.Redis.Port = "6379"
	cfg.Redis.DB = 0

	cfg.Elastic.Host = "localhost"
	cfg.Elastic.Port = "9200"

	cfg.Cache.DefaultTtl = time.Hour
	cfg.Cache.CountersTtl = time.Hour

	cfg.Tokens.AccessTTL = time.Hour
	cfg.Tokens.RefreshTTL = time.Hour
	cfg.Tokens.RecoveryTTL = time.Hour

	cfg.JWT.PrivateKey = nil
	err = cfg.Validate()
	assert.Error(t, err)
}
