package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig_NotNil(t *testing.T) {
	config, err := NewServerConfig()
	require.NoError(t, err)

	assert.NotNil(t, config)
}

func TestNewAgentConfig_NotNil(t *testing.T) {
	config, err := NewAgentConfig()
	require.NoError(t, err)

	assert.NotNil(t, config)
}

func TestServerConfig_Interface(t *testing.T) {
	config, err := NewServerConfig()
	require.NoError(t, err)

	// Test that config implements ServerConfig interface
	var serverConfig ServerConfig = config
	assert.NotNil(t, serverConfig)
}

func TestAgentConfig_Interface(t *testing.T) {
	config, err := NewAgentConfig()
	require.NoError(t, err)

	// Test that config implements AgentClientConfig interface
	var clientConfig AgentClientConfig = config
	assert.NotNil(t, clientConfig)
}

func TestGetProjectRoot_ValidProject(t *testing.T) {
	root, err := GetProjectRoot()

	assert.NoError(t, err)
	assert.NotEmpty(t, root)

	goModPath := filepath.Join(root, "go.mod")
	_, err = os.Stat(goModPath)
	assert.NoError(t, err)
}

func TestGetProjectRoot_FromTempDir(t *testing.T) {
	// Create temporary directory without go.mod
	tempDir := t.TempDir()

	// Save current directory
	originalWd, err := os.Getwd()
	assert.NoError(t, err)

	// Change to temporary directory
	err = os.Chdir(tempDir)
	assert.NoError(t, err)

	// Restore original directory after test
	defer func() {
		os.Chdir(originalWd)
	}()

	// Try to find project root - should return error
	_, err = GetProjectRoot()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "project root not found")
}

func TestGetEncFilePath_NotEmpty(t *testing.T) {
	path := getEncFilePath()

	assert.NotEmpty(t, path)
	assert.Contains(t, path, ".env")
}

func TestGetEncFilePath_WithMockedGetEnvPath(t *testing.T) {
	// Save original function
	originalGetEnvPath := getEnvPath

	// Mock the function
	getEnvPath = func() string {
		return "/mock/path/.env"
	}

	// Restore after test
	defer func() {
		getEnvPath = originalGetEnvPath
	}()

	path := getEnvPath()
	assert.Equal(t, "/mock/path/.env", path)
}

func TestNewServerConfig_WithMockedEnv(t *testing.T) {
	// Mock getEnvPath to avoid file issues
	originalGetEnvPath := getEnvPath
	getEnvPath = func() string {
		return "/nonexistent/.env" // File doesn't exist, but this shouldn't break config creation
	}
	defer func() {
		getEnvPath = originalGetEnvPath
	}()

	config, err := NewServerConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestNewAgentConfig_WithMockedEnv(t *testing.T) {
	// Mock getEnvPath to avoid file issues
	originalGetEnvPath := getEnvPath
	getEnvPath = func() string {
		return "/nonexistent/.env" // File doesn't exist, but this shouldn't break config creation
	}
	defer func() {
		getEnvPath = originalGetEnvPath
	}()

	config, err := NewAgentConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestConfig_Methods(t *testing.T) {
	config := &Config{}

	// Test setting and getting via methods
	config.Addr = "localhost:8080"
	config.SecretKey = "secret"

	assert.Equal(t, "localhost:8080", config.GetAddress())
	assert.Equal(t, "secret", config.GetSecretKey())

	// Test agent config methods
	err := config.SetMasterPassword("master")
	assert.NoError(t, err)
	masterPass, err := config.GetMasterPassword()
	assert.NoError(t, err)
	assert.Equal(t, "master", masterPass)

	err = config.SetMasterKey([]byte("key"))
	assert.NoError(t, err)
	masterKey, err := config.GetMasterKey()
	assert.NoError(t, err)
	assert.Equal(t, []byte("key"), masterKey)

	err = config.SetSalt([]byte("salt"))
	assert.NoError(t, err)
	salt, err := config.GetSalt()
	assert.NoError(t, err)
	assert.Equal(t, []byte("salt"), salt)

	// Test server config methods
	config.DBConnStr = "postgres://..."
	assert.Equal(t, "postgres://...", config.GetConnectionString())

	err = config.SetPublicKeyPEM([]byte("pem"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("pem"), config.GetPublicKeyPEM())
}

func TestConfig_ZeroValues(t *testing.T) {
	config := &Config{}

	// Test zero values via methods
	assert.Empty(t, config.GetAddress())
	assert.Empty(t, config.GetSecretKey())
	assert.Empty(t, config.GetConnectionString())
	assert.Nil(t, config.GetPrivateKey())
	assert.Nil(t, config.GetPublicKeyPEM())

	// Test methods that should return errors for empty values
	_, err := config.GetMasterPassword()
	assert.Error(t, err)

	_, err = config.GetSalt()
	assert.Error(t, err)
}
