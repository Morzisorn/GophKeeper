package controllers

import (
	"context"
	"testing"

	"gophkeeper/config"
	pb "gophkeeper/internal/protos/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCryptoController(t *testing.T) {
	cnfg, err := config.NewServerConfig()
	require.NoError(t, err)

	controller := NewCryptoController(cnfg)
	assert.NotNil(t, controller)
}

func TestCryptoController_GetPublicKeyPEM(t *testing.T) {
	cnfg, err := config.NewServerConfig()
	require.NoError(t, err)

	controller := NewCryptoController(cnfg)
	request := &pb.GetPublicKeyPEMRequest{}

	response, err := controller.GetPublicKeyPEM(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}
