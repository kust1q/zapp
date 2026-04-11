package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kust1q/Zapp/main/internal/config"
)

func NewPostgresConnect(cfg *config.PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
