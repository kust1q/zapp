package config

import (
	"crypto/rsa"
	"time"
)

type AuthServiceConfig struct {
	PrivateKey  *rsa.PrivateKey
	PublicKey   *rsa.PublicKey
	AccessTTL   time.Duration
	RefreshTTL  time.Duration
	RecoveryTTL time.Duration
}
