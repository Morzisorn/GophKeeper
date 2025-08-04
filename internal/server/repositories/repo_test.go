package repositories

import (
	"testing"

	"gophkeeper/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDatabaseConfig struct {
	dbConnStr string
}

func (m *mockDatabaseConfig) GetConnectionString() string {
	return m.dbConnStr
}

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name      string
		dbConnStr string
		wantErr   bool
	}{
		{
			name:      "invalid connection string",
			dbConnStr: "invalid_connection_string",
			wantErr:   true,
		},
		{
			name:      "empty connection string",
			dbConnStr: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &mockDatabaseConfig{dbConnStr: tt.dbConnStr}
			storage, err := NewStorage(cfg)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, storage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, storage)
			}
		})
	}
}

func TestNewStorage_WithRealConfig(t *testing.T) {
	cfg, err := config.NewServerConfig()
	require.NoError(t, err)

	// Test with real config - this should succeed in creating the storage object
	// but may fail when actually connecting to database
	storage, err := NewStorage(cfg)
	// With the default config, this should succeed (no error expected)
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}
