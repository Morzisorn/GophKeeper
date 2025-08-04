package client

import (
	"github.com/stretchr/testify/require"
	"gophkeeper/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGRPCClient_NotNil(t *testing.T) {
	cnfg, err := config.NewAgentConfig()
	require.NoError(t, err)
	cnfg.Addr = "localhost:8080"

	client, err := NewGRPCClient(cnfg)
	require.NoError(t, err)

	assert.NotNil(t, client)
}
