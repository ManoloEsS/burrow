package paths

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const appName = "burrow"

func GetConfigPath() string {
	return filepath.Join(xdg.ConfigHome, appName, "config.yaml")
}

func GetDatabasePath() string {
	return filepath.Join(xdg.DataHome, appName, "burrow.db")
}

func GetLogPath() string {
	return filepath.Join(xdg.StateHome, appName, "burrow_log")
}

func GetCachePath() string {
	return filepath.Join(xdg.CacheHome, appName)
}

func GetServerCachePath() string {
	return filepath.Join(GetCachePath(), "servers")
}

func GetConfigDir() string {
	return filepath.Dir(GetConfigPath())
}

func GetDataDir() string {
	return filepath.Dir(GetDatabasePath())
}

func GetStateDir() string {
	return filepath.Dir(GetLogPath())
}

func EnsureDirectories() error {
	dirs := []string{
		GetConfigDir(),
		GetDataDir(),
		GetStateDir(),
		GetCachePath(),
		GetServerCachePath(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func SearchConfigFile() (string, error) {
	return xdg.SearchConfigFile(filepath.Join(appName, "config.yaml"))
}

func ConfigFileExists() bool {
	configPath := GetConfigPath()
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}
