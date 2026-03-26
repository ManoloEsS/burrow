package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
)

func TestLoad_WithDefaults(t *testing.T) {
	withIsolatedXDG(t, func() {
		clearEnvVars()

		cfg, err := Load()
		if err != nil {
			t.Logf("Load error: %v", err)
		}
		assert.NoError(t, err)
		assert.NotNil(t, cfg)

		assert.Equal(t, "8080", cfg.App.DefaultPort)
		assert.Equal(t, GetDatabasePath(), cfg.Database.Path)
		assert.Equal(t, GetConfigPath(), cfg.Paths.ConfigFile)
		assert.Equal(t, GetLogPath(), cfg.Paths.LogFile)

		expectedConnectionString := fmt.Sprintf(
			"file:%s?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL",
			GetDatabasePath(),
		)
		assert.Equal(t, expectedConnectionString, cfg.Database.ConnectionString)
	})
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	withIsolatedXDG(t, func() {
		clearEnvVars()

		os.Setenv("BURROW_DEFAULT_PORT", "3000")
		os.Setenv("BURROW_DB_FILE", "/tmp/test.db")

		t.Cleanup(clearEnvVars)

		cfg, err := Load()
		assert.NoError(t, err)

		assert.Equal(t, "3000", cfg.App.DefaultPort)
		assert.Equal(t, "/tmp/test.db", cfg.Database.Path)

		expectedConnectionString := "file:/tmp/test.db?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL"
		assert.Equal(t, expectedConnectionString, cfg.Database.ConnectionString)
	})
}

func TestLoad_WithConfigFile(t *testing.T) {
	withIsolatedXDG(t, func() {
		clearEnvVars()

		writeConfigFile(t, `app:
	  default_port: "9000"
database:
	  path: "/custom/path/db.sqlite"
`)

		cfg, err := Load()
		assert.NoError(t, err)

		assert.Equal(t, "9000", cfg.App.DefaultPort)
		assert.Equal(t, "/custom/path/db.sqlite", cfg.Database.Path)
	})
}

func TestLoad_EnvironmentOverridesConfig(t *testing.T) {
	withIsolatedXDG(t, func() {
		clearEnvVars()

		writeConfigFile(t, `app:
	  default_port: "9000"
database:
	  path: "/custom/path/db.sqlite"
`)

		os.Setenv("BURROW_DEFAULT_PORT", "5000")
		defer clearEnvVars()

		cfg, err := Load()
		assert.NoError(t, err)

		assert.Equal(t, "5000", cfg.App.DefaultPort)
		assert.Equal(t, "/custom/path/db.sqlite", cfg.Database.Path)
	})
}

func TestEnsureDirectories(t *testing.T) {
	withIsolatedXDG(t, func() {
		assert.NoError(t, EnsureDirectories())
	})

	t.Run("mkdir failure surfaces error", func(t *testing.T) {
		withIsolatedXDG(t, func() {
			orig := mkdirAll
			mkdirAll = func(string, os.FileMode) error {
				return errors.New("forced failure")
			}
			defer func() { mkdirAll = orig }()

			err := EnsureDirectories()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "forced failure")
		})
	})
}

func TestValidate_MissingDatabasePath(t *testing.T) {
	cfg := &Config{
		App: AppConfig{DefaultPort: "8080"},
		Database: DatabaseConfig{
			Path: "",
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
			Path: "/path/to/db.sqlite",
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
	envVars := []string{"BURROW_DEFAULT_PORT", "BURROW_DB_FILE", "XDG_CONFIG_HOME", "XDG_DATA_HOME", "XDG_STATE_HOME", "XDG_CACHE_HOME"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}
}

func withIsolatedXDG(t *testing.T, fn func()) {
	t.Helper()
	clearEnvVars()
	tmpBase := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpBase, "config"))
	os.Setenv("XDG_DATA_HOME", filepath.Join(tmpBase, "data"))
	os.Setenv("XDG_STATE_HOME", filepath.Join(tmpBase, "state"))
	os.Setenv("XDG_CACHE_HOME", filepath.Join(tmpBase, "cache"))

	// xdg package caches env vars at init; reset before running logic
	xdg.Reload()

	t.Cleanup(func() {
		clearEnvVars()
		xdg.Reload()
	})

	fn()
}

func writeConfigFile(t *testing.T, contents string) {
	t.Helper()
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)
	assert.NoError(t, os.MkdirAll(configDir, 0o755))
	data := strings.TrimSpace(strings.ReplaceAll(contents, "\t", "")) + "\n"
	assert.NoError(t, os.WriteFile(configPath, []byte(data), 0o644))
	t.Cleanup(func() {
		_ = os.Remove(configPath)
	})
}
