package config

import (
	"fmt"
	"os"
)

type Config struct {
	DbPath          string
	DefaultPort     string
	DbMigrationsDir string
	DbString        string
}

func LoadFromEnv() *Config {
	cfg := &Config{
		DbPath:          os.Getenv("DB_FILE"),
		DefaultPort:     os.Getenv("DEFAULT_PORT"),
		DbMigrationsDir: os.Getenv("GOOSE_MIGRATIONS_DIR"),
		DbString:        os.Getenv("GOOSE_DBSTRING"),
	}

	if cfg.DbMigrationsDir == "" {
		cfg.DbMigrationsDir = "sql/migrations"
	}

	if cfg.DbString == "" && cfg.DbPath != "" {
		cfg.DbString = fmt.Sprintf(
			"file:%s?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL",
			cfg.DbPath,
		)
	}

	return cfg
}

func (c *Config) Validate() error {
	if c.DbPath == "" {
		return fmt.Errorf("DB_FILE environment variable must be set")
	}
	if c.DefaultPort == "" {
		return fmt.Errorf("DEFAULT_PORT environment variable must be set")
	}
	if c.DbString == "" {
		return fmt.Errorf("GOOSE_DBSTRING environment variable must be set")
	}
	return nil
}
