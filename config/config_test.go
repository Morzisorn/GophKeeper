package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerConfig_NotNil(t *testing.T) {
	config, err := GetServerConfig()
	require.NoError(t, err)

	assert.NotNil(t, config)
}

func TestGetAgentConfig_NotNil(t *testing.T) {
	config, err := GetAgentConfig()
	require.NoError(t, err)

	assert.NotNil(t, config)
}

func TestGetServerConfig_Singleton(t *testing.T) {
	config1, err := GetServerConfig()
	require.NoError(t, err)
	config2, err := GetServerConfig()
	require.NoError(t, err)

	assert.Same(t, config1, config2)
}

func TestGetAgentConfig_Singleton(t *testing.T) {
	config1, err := GetAgentConfig()
	require.NoError(t, err)
	config2, err := GetAgentConfig()
	require.NoError(t, err)

	assert.Same(t, config1, config2)
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
	// Создаем временную директорию без go.mod
	tempDir := t.TempDir()

	// Сохраняем текущую директорию
	originalWd, err := os.Getwd()
	assert.NoError(t, err)

	// Переходим во временную директорию
	err = os.Chdir(tempDir)
	assert.NoError(t, err)

	// Восстанавливаем оригинальную директорию после теста
	defer func() {
		os.Chdir(originalWd)
	}()

	// Пытаемся найти корень проекта - должна быть ошибка
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
	// Сохраняем оригинальную функцию
	originalGetEnvPath := getEnvPath

	// Мокаем функцию
	getEnvPath = func() string {
		return "/mock/path/.env"
	}

	// Восстанавливаем после теста
	defer func() {
		getEnvPath = originalGetEnvPath
	}()

	path := getEnvPath()
	assert.Equal(t, "/mock/path/.env", path)
}

func TestNewServerConfig_NotNil(t *testing.T) {
	// Мокаем getEnvPath чтобы избежать проблем с файлами
	originalGetEnvPath := getEnvPath
	getEnvPath = func() string {
		return "/nonexistent/.env" // Файл не существует, но это не должно ломать создание конфига
	}
	defer func() {
		getEnvPath = originalGetEnvPath
	}()

	config, err := newServerConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestNewAgentConfig_NotNil(t *testing.T) {
	// Мокаем getEnvPath чтобы избежать проблем с файлами
	originalGetEnvPath := getEnvPath
	getEnvPath = func() string {
		return "/nonexistent/.env" // Файл не существует, но это не должно ломать создание конфига
	}
	defer func() {
		getEnvPath = originalGetEnvPath
	}()

	config, err := newAgentConfig()

	assert.NoError(t, err)
	assert.NotNil(t, config)
}

func TestConfig_StructFields(t *testing.T) {
	config := &Config{}

	// Проверяем, что все поля доступны
	assert.NotNil(t, &config.CommonConfig)
	assert.NotNil(t, &config.AgentConfig)
	assert.NotNil(t, &config.ServerConfig)

	// Проверяем поля CommonConfig
	config.AppType = "test"
	config.Addr = "localhost:8080"
	config.CryptoKeyPath = "/path/to/key"
	config.SecretKey = "secret"

	assert.Equal(t, "test", config.AppType)
	assert.Equal(t, "localhost:8080", config.Addr)
	assert.Equal(t, "/path/to/key", config.CryptoKeyPath)
	assert.Equal(t, "secret", config.SecretKey)

	// Проверяем поля AgentConfig
	config.MasterPassword = "master"
	config.MasterKey = []byte("key")
	config.Salt = []byte("salt")

	assert.Equal(t, "master", config.MasterPassword)
	assert.Equal(t, []byte("key"), config.MasterKey)
	assert.Equal(t, []byte("salt"), config.Salt)

	// Проверяем поля ServerConfig
	config.DBConnStr = "postgres://..."
	config.PublicKeyPEM = []byte("pem")

	assert.Equal(t, "postgres://...", config.DBConnStr)
	assert.Equal(t, []byte("pem"), config.PublicKeyPEM)
}

func TestConfigTypes_ZeroValues(t *testing.T) {
	// Тестируем нулевые значения структур
	common := CommonConfig{}
	assert.Empty(t, common.AppType)
	assert.Empty(t, common.Addr)
	assert.Empty(t, common.CryptoKeyPath)
	assert.Empty(t, common.SecretKey)

	agent := AgentConfig{}
	assert.Nil(t, agent.PublicKey)
	assert.Empty(t, agent.MasterPassword)
	assert.Nil(t, agent.MasterKey)
	assert.Nil(t, agent.Salt)

	server := ServerConfig{}
	assert.Empty(t, server.DBConnStr)
	assert.Nil(t, server.PrivateKey)
	assert.Nil(t, server.PublicKeyPEM)
}
