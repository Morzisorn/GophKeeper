package controllers

import (
	"context"
	"testing"

	pb "gophkeeper/internal/protos/crypto"

	"github.com/stretchr/testify/assert"
)

func TestNewCryptoController(t *testing.T) {
	controller := NewCryptoController()
	assert.NotNil(t, controller)
}

func TestCryptoController_GetPublicKeyPEM(t *testing.T) {
	controller := NewCryptoController()
	request := &pb.GetPublicKeyPEMRequest{}

	response, err := controller.GetPublicKeyPEM(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}
