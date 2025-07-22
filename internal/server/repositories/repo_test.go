package repositories

import (
	"testing"

	"gophkeeper/config"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "invalid connection string",
			cfg: &config.Config{
				ServerConfig: config.ServerConfig{
					DBConnStr: "invalid_connection_string",
				},
			},
			wantErr: true,
		},
		{
			name: "empty connection string",
			cfg: &config.Config{
				ServerConfig: config.ServerConfig{
					DBConnStr: "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage, err := NewStorage(tt.cfg)
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
