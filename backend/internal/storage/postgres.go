package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB() (*sql.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "launchpad")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "launchpad")
	sslmode := getEnv("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			address VARCHAR(42) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS tokens (
			id SERIAL PRIMARY KEY,
			address VARCHAR(42) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			symbol VARCHAR(10) NOT NULL,
			total_supply VARCHAR(255) NOT NULL,
			creator_address VARCHAR(42) NOT NULL,
			tx_hash VARCHAR(66) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS presales (
			id SERIAL PRIMARY KEY,
			address VARCHAR(42) UNIQUE NOT NULL,
			token_address VARCHAR(42) NOT NULL,
			creator_address VARCHAR(42) NOT NULL,
			rate VARCHAR(255) NOT NULL,
			soft_cap VARCHAR(255) NOT NULL,
			hard_cap VARCHAR(255) NOT NULL,
			deadline TIMESTAMP NOT NULL,
			tx_hash VARCHAR(66) NOT NULL,
			active BOOLEAN DEFAULT true,
			finalized BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS presale_participations (
			id SERIAL PRIMARY KEY,
			presale_id INTEGER NOT NULL REFERENCES presales(id),
			participant_address VARCHAR(42) NOT NULL,
			amount_eth VARCHAR(255) NOT NULL,
			amount_tokens VARCHAR(255) NOT NULL,
			tx_hash VARCHAR(66) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_creator ON tokens(creator_address)`,
		`CREATE INDEX IF NOT EXISTS idx_presales_creator ON presales(creator_address)`,
		`CREATE INDEX IF NOT EXISTS idx_presales_token ON presales(token_address)`,
		`CREATE INDEX IF NOT EXISTS idx_participations_presale ON presale_participations(presale_id)`,
		`CREATE INDEX IF NOT EXISTS idx_participations_participant ON presale_participations(participant_address)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to run migration %d: %w", i+1, err)
		}
	}

	return nil
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}