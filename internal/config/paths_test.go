package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigPath(t *testing.T) {
	path := GetConfigPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "burrow")
}

func TestGetDatabasePath(t *testing.T) {
	path := GetDatabasePath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "burrow.db")
}

func TestGetLogPath(t *testing.T) {
	path := GetLogPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "burrow_log")
}

func TestGetServerCachePath(t *testing.T) {
	path := GetServerCachePath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "servers")
}

func TestSearchConfigFile(t *testing.T) {
	withTempXDG(t, func() {
		cfgPath := GetConfigPath()
		assert.NoError(t, os.MkdirAll(filepath.Dir(cfgPath), 0o755))
		assert.NoError(t, os.WriteFile(cfgPath, []byte("foo: bar\n"), 0o644))

		foundPath, err := SearchConfigFile()
		assert.NoError(t, err)
		assert.Equal(t, cfgPath, foundPath)
	})

	withTempXDG(t, func() {
		_, err := SearchConfigFile()
		assert.Error(t, err)
	})
}

func TestConfigFileExists(t *testing.T) {
	withTempXDG(t, func() {
		cfgPath := GetConfigPath()
		assert.NoError(t, os.MkdirAll(filepath.Dir(cfgPath), 0o755))
		assert.NoError(t, os.WriteFile(cfgPath, []byte("foo"), 0o644))

		assert.True(t, ConfigFileExists())
	})

	withTempXDG(t, func() {
		assert.False(t, ConfigFileExists())
	})
}

func TestEnsureDirectoriesCreatesMissing(t *testing.T) {
	withTempXDG(t, func() {
		assert.NoError(t, EnsureDirectories())

		for _, dir := range []string{GetConfigDir(), GetCachePath(), GetServerCachePath()} {
			info, err := os.Stat(dir)
			assert.NoError(t, err)
			assert.True(t, info.IsDir())
		}
	})
}

func withTempXDG(t *testing.T, fn func()) {
	t.Helper()
	tmp := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "config"))
	os.Setenv("XDG_DATA_HOME", filepath.Join(tmp, "data"))
	os.Setenv("XDG_STATE_HOME", filepath.Join(tmp, "state"))
	os.Setenv("XDG_CACHE_HOME", filepath.Join(tmp, "cache"))
	xdg.Reload()

	t.Cleanup(func() {
		os.Unsetenv("XDG_CONFIG_HOME")
		os.Unsetenv("XDG_DATA_HOME")
		os.Unsetenv("XDG_STATE_HOME")
		os.Unsetenv("XDG_CACHE_HOME")
		xdg.Reload()
	})

	fn()
}
