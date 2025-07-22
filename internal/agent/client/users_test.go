package client

import (
	"context"
	"gophkeeper/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGRPCClient_SignUpUser_NilClient(t *testing.T) {
	var client *GRPCClient = nil
	
	user := &models.User{
		Login:    "test-login",
		Password: []byte("test-password"),
	}
	
	assert.Panics(t, func() {
		client.SignUpUser(context.Background(), user)
	})
}

func TestGRPCClient_SignUpUser_NilUserClient(t *testing.T) {
	client := &GRPCClient{
		User: nil,
	}
	
	user := &models.User{
		Login:    "test-login",
		Password: []byte("test-password"),
	}
	
	assert.Panics(t, func() {
		client.SignUpUser(context.Background(), user)
	})
}

func TestGRPCClient_SignInUser_NilClient(t *testing.T) {
	var client *GRPCClient = nil
	
	user := &models.User{
		Login:    "test-login",
		Password: []byte("test-password"),
	}
	
	assert.Panics(t, func() {
		client.SignInUser(context.Background(), user)
	})
}

func TestGRPCClient_SignInUser_NilUserClient(t *testing.T) {
	client := &GRPCClient{
		User: nil,
	}
	
	user := &models.User{
		Login:    "test-login",
		Password: []byte("test-password"),
	}
	
	assert.Panics(t, func() {
		client.SignInUser(context.Background(), user)
	})
}

func TestGRPCClient_SetJWTToken_NilClient(t *testing.T) {
	var client *GRPCClient = nil
	
	assert.Panics(t, func() {
		client.SetJWTToken("test-token")
	})
}

func TestGRPCClient_GetJWTToken_NilClient(t *testing.T) {
	var client *GRPCClient = nil
	
	assert.Panics(t, func() {
		client.GetJWTToken()
	})
}

func TestGRPCClient_SetJWTToken_ValidClient(t *testing.T) {
	client := &GRPCClient{}
	
	token := "test-jwt-token"
	client.SetJWTToken(token)
	
	assert.Equal(t, token, client.token)
}

func TestGRPCClient_GetJWTToken_ValidClient(t *testing.T) {
	client := &GRPCClient{
		token: "test-jwt-token",
	}
	
	result := client.GetJWTToken()
	
	assert.Equal(t, "test-jwt-token", result)
}