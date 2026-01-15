package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromEnv(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		expectedConfig *Config
		description    string
	}{
		{
			name: "All environment variables set",
			envVars: map[string]string{
				"DB_FILE":              "/path/to/db.sqlite",
				"DEFAULT_PORT":         "8080",
				"GOOSE_MIGRATIONS_DIR": "custom/migrations",
				"GOOSE_DBSTRING":       "custom:connection:string",
			},
			expectedConfig: &Config{
				DbPath:          "/path/to/db.sqlite",
				DefaultPort:     "8080",
				DbMigrationsDir: "custom/migrations",
				DbString:        "custom:connection:string",
			},
			description: "Should load all environment variables correctly",
		},
		{
			name: "Default migrations directory",
			envVars: map[string]string{
				"DB_FILE":              "/path/to/db.sqlite",
				"DEFAULT_PORT":         "8080",
				"GOOSE_DBSTRING":       "file:/path/to/db.sqlite?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL",
				"GOOSE_MIGRATIONS_DIR": "",
			},
			expectedConfig: &Config{
				DbPath:          "/path/to/db.sqlite",
				DefaultPort:     "8080",
				DbMigrationsDir: "sql/migrations",
				DbString:        "file:/path/to/db.sqlite?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL",
			},
			description: "Should use default migrations directory when not set",
		},
		{
			name: "Auto-generate DB string from DB_PATH",
			envVars: map[string]string{
				"DB_FILE":              "/path/to/db.sqlite",
				"DEFAULT_PORT":         "8080",
				"GOOSE_DBSTRING":       "",
				"GOOSE_MIGRATIONS_DIR": "sql/migrations",
			},
			expectedConfig: &Config{
				DbPath:          "/path/to/db.sqlite",
				DefaultPort:     "8080",
				DbMigrationsDir: "sql/migrations",
				DbString:        "file:/path/to/db.sqlite?cache=shared&mode=rwc&_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL",
			},
			description: "Should auto-generate DB string when DB_PATH is set but DB_STRING is not",
		},
		{
			name: "Empty environment variables",
			envVars: map[string]string{
				"DB_FILE":              "",
				"DEFAULT_PORT":         "",
				"GOOSE_MIGRATIONS_DIR": "",
				"GOOSE_DBSTRING":       "",
			},
			expectedConfig: &Config{
				DbPath:          "",
				DefaultPort:     "",
				DbMigrationsDir: "sql/migrations",
				DbString:        "",
			},
			description: "Should handle empty environment variables",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key := range tt.envVars {
				os.Unsetenv(key)
			}

			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			defer func() {
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			config := LoadFromEnv()

			assert.Equal(t, tt.expectedConfig.DbPath, config.DbPath, "DB_PATH should match")
			assert.Equal(t, tt.expectedConfig.DefaultPort, config.DefaultPort, "DEFAULT_PORT should match")
			assert.Equal(t, tt.expectedConfig.DbMigrationsDir, config.DbMigrationsDir, "GOOSE_MIGRATIONS_DIR should match")
			assert.Equal(t, tt.expectedConfig.DbString, config.DbString, "GOOSE_DBSTRING should match")
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
		errorField  string
	}{
		{
			name: "Valid configuration",
			config: &Config{
				DbPath:      "/path/to/db.sqlite",
				DefaultPort: "8080",
				DbString:    "file:/path/to/db.sqlite?cache=shared",
			},
			expectError: false,
			errorField:  "",
		},
		{
			name: "Missing DB_FILE",
			config: &Config{
				DbPath:      "",
				DefaultPort: "8080",
				DbString:    "file:/path/to/db.sqlite?cache=shared",
			},
			expectError: true,
			errorField:  "DB_FILE",
		},
		{
			name: "Missing DEFAULT_PORT",
			config: &Config{
				DbPath:      "/path/to/db.sqlite",
				DefaultPort: "",
				DbString:    "file:/path/to/db.sqlite?cache=shared",
			},
			expectError: true,
			errorField:  "DEFAULT_PORT",
		},
		{
			name: "Missing GOOSE_DBSTRING",
			config: &Config{
				DbPath:      "/path/to/db.sqlite",
				DefaultPort: "8080",
				DbString:    "",
			},
			expectError: true,
			errorField:  "GOOSE_DBSTRING",
		},
		{
			name: "All fields missing",
			config: &Config{
				DbPath:      "",
				DefaultPort: "",
				DbString:    "",
			},
			expectError: true,
			errorField:  "DB_FILE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err, "Expected validation error")
				assert.Contains(t, err.Error(), tt.errorField, "Error message should mention missing field")
			} else {
				assert.NoError(t, err, "Expected no validation error")
			}
		})
	}
}
