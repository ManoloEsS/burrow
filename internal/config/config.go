package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Paths    PathsConfig    `yaml:"-"`
}

type AppConfig struct {
	DefaultPort string `yaml:"default_port"`
}

type DatabaseConfig struct {
	Path             string `yaml:"path"`
	MigrationsDir    string `yaml:"migrations_dir"`
	ConnectionString string `yaml:"-"`
}

type PathsConfig struct {
	ConfigFile string `yaml:"-"`
	LogFile    string `yaml:"-"`
}

func Load() (*Config, error) {
	cfg := &Config{}

	applyDefaults(cfg)

	err := loadFromFile(cfg)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("config file error: %w", err)
		}
	}

	loadFromEnv(cfg)

	if err := resolvePaths(cfg); err != nil {
		return nil, fmt.Errorf("path resolution error: %w", err)
	}

	cfg.Database.ConnectionString = generateDbString(cfg.Database.Path)

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func applyDefaults(cfg *Config) {
	cfg.App.DefaultPort = "8080"
	cfg.Database.MigrationsDir = ""
	cfg.Database.Path = ""
}

func loadFromFile(cfg *Config) error {
	if !ConfigFileExists() {
		return os.ErrNotExist
	}

	configPath, err := SearchConfigFile()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}

func loadFromEnv(cfg *Config) {
	if port := os.Getenv("DEFAULT_PORT"); port != "" {
		cfg.App.DefaultPort = port
	}

	if dbPath := os.Getenv("DB_FILE"); dbPath != "" {
		cfg.Database.Path = dbPath
	}

	// GOOSE_MIGRATIONS_DIR is deprecated but still supported for development
	if migrationsDir := os.Getenv("GOOSE_MIGRATIONS_DIR"); migrationsDir != "" {
		log.Printf("Warning: GOOSE_MIGRATIONS_DIR is deprecated. Migrations are now embedded.")
		cfg.Database.MigrationsDir = migrationsDir
	}

}

func resolvePaths(cfg *Config) error {
	if cfg.Database.Path == "" {
		cfg.Database.Path = GetDatabasePath()
	}

	cfg.Paths.ConfigFile = GetConfigPath()
	cfg.Paths.LogFile = GetLogPath()

	if err := EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure directories: %w", err)
	}

	return nil
}

func generateDbString(dbPath string) string {
	return fmt.Sprintf(
		"file:%s?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL",
		dbPath,
	)
}

func validate(cfg *Config) error {
	if cfg.Database.Path == "" {
		return fmt.Errorf("database path cannot be empty")
	}

	if cfg.App.DefaultPort == "" {
		return fmt.Errorf("default port cannot be empty")
	}

	// migrations_dir is now optional (embedded migrations are default)
	if cfg.Database.MigrationsDir != "" {
		if _, err := os.Stat(cfg.Database.MigrationsDir); os.IsNotExist(err) {
			return fmt.Errorf("migrations directory does not exist: %s", cfg.Database.MigrationsDir)
		}
	}

	return nil
}
