package paths

import (
	"testing"

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

func TestGetCachePath(t *testing.T) {
	path := GetCachePath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "burrow")
}

func TestGetServerCachePath(t *testing.T) {
	path := GetServerCachePath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "servers")
}
