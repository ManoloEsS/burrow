package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/ManoloEsS/burrow/internal/pkg/paths"
	"github.com/stretchr/testify/assert"
)

func TestLoad_WithDefaults(t *testing.T) {
	clearEnvVars()

	cfg, err := Load()
	if err != nil {
		t.Logf("Load error: %v", err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "8080", cfg.App.DefaultPort)
	assert.Equal(t, "sql/migrations", cfg.Database.MigrationsDir)
	assert.Equal(t, paths.GetDatabasePath(), cfg.Database.Path)
	assert.Equal(t, paths.GetConfigPath(), cfg.Paths.ConfigFile)
	assert.Equal(t, paths.GetLogPath(), cfg.Paths.LogFile)

	expectedConnectionString := fmt.Sprintf(
		"file:%s?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL",
		paths.GetDatabasePath(),
	)
	assert.Equal(t, expectedConnectionString, cfg.Database.ConnectionString)
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	clearEnvVars()

	os.Setenv("DEFAULT_PORT", "3000")
	os.Setenv("DB_FILE", "/tmp/test.db")
	os.Setenv("GOOSE_MIGRATIONS_DIR", "custom/migrations")

	defer clearEnvVars()

	cfg, err := Load()
	assert.NoError(t, err)

	assert.Equal(t, "3000", cfg.App.DefaultPort)
	assert.Equal(t, "/tmp/test.db", cfg.Database.Path)
	assert.Equal(t, "custom/migrations", cfg.Database.MigrationsDir)

	expectedConnectionString := fmt.Sprintf(
		"file:/tmp/test.db?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL",
	)
	assert.Equal(t, expectedConnectionString, cfg.Database.ConnectionString)
}

func TestLoad_WithConfigFile(t *testing.T) {
	clearEnvVars()

	configPath := paths.GetConfigPath()
	configDir := paths.GetConfigDir()

	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(configDir)

	configContent := `app:
  default_port: "9000"
database:
  path: "/custom/path/db.sqlite"
  migrations_dir: "custom/migrations"
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(configPath)

	cfg, err := Load()
	assert.NoError(t, err)

	assert.Equal(t, "9000", cfg.App.DefaultPort)
	assert.Equal(t, "/custom/path/db.sqlite", cfg.Database.Path)
	assert.Equal(t, "custom/migrations", cfg.Database.MigrationsDir)
}

func TestLoad_EnvironmentOverridesConfig(t *testing.T) {
	clearEnvVars()

	configPath := paths.GetConfigPath()
	configDir := paths.GetConfigDir()

	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)
	defer os.RemoveAll(configDir)

	configContent := `app:
  default_port: "9000"
database:
  path: "/custom/path/db.sqlite"
`
	err = os.WriteFile(configPath, []byte(configContent), 0644)
	assert.NoError(t, err)
	defer os.Remove(configPath)

	os.Setenv("DEFAULT_PORT", "5000")
	defer clearEnvVars()

	cfg, err := Load()
	assert.NoError(t, err)

	assert.Equal(t, "5000", cfg.App.DefaultPort)
	assert.Equal(t, "/custom/path/db.sqlite", cfg.Database.Path)
}

func TestValidate_MissingDatabasePath(t *testing.T) {
	cfg := &Config{
		App: AppConfig{DefaultPort: "8080"},
		Database: DatabaseConfig{
			Path:          "",
			MigrationsDir: "sql/migrations",
		},
	}

	err := validate(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database path cannot be empty")
}

func TestValidate_MissingDefaultPort(t *testing.T) {
	cfg := &Config{
		App: AppConfig{DefaultPort: ""},
		Database: DatabaseConfig{
			Path:          "/path/to/db.sqlite",
			MigrationsDir: "sql/migrations",
		},
	}

	err := validate(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "default port cannot be empty")
}

func TestGenerateDbString(t *testing.T) {
	dbPath := "/path/to/test.db"
	expected := "file:/path/to/test.db?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL"

	result := generateDbString(dbPath)
	assert.Equal(t, expected, result)
}

func clearEnvVars() {
	envVars := []string{"DEFAULT_PORT", "DB_FILE", "GOOSE_MIGRATIONS_DIR"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}
}
