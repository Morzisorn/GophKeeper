package client

import (
	"gophkeeper/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGRPCClient_NotNil(t *testing.T) {
	cfg := config.GetServerConfig()
	cfg.Addr = "localhost:8080"
	
	client := NewGRPCClient(cfg)
	
	assert.NotNil(t, client)
}