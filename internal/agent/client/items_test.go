package client

import (
	"context"
	"gophkeeper/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGRPCClient_AddItem_NilClient(t *testing.T) {
	var client *GRPCClient = nil

	item := &models.EncryptedItem{}

	assert.Panics(t, func() {
		client.AddItem(context.Background(), item)
	})
}

func TestGRPCClient_AddItem_NilItemClient(t *testing.T) {
	client := &GRPCClient{
		Item: nil,
	}

	item := &models.EncryptedItem{}

	assert.Panics(t, func() {
		client.AddItem(context.Background(), item)
	})
}

func TestGRPCClient_EditItem_NilClient(t *testing.T) {
	var client *GRPCClient = nil

	item := &models.EncryptedItem{}

	assert.Panics(t, func() {
		client.EditItem(context.Background(), item)
	})
}

func TestGRPCClient_EditItem_NilItemClient(t *testing.T) {
	client := &GRPCClient{
		Item: nil,
	}

	item := &models.EncryptedItem{}

	assert.Panics(t, func() {
		client.EditItem(context.Background(), item)
	})
}

func TestGRPCClient_DeleteItem_NilClient(t *testing.T) {
	var client *GRPCClient = nil

	var itemID [16]byte

	assert.Panics(t, func() {
		client.DeleteItem(context.Background(), "test-login", itemID)
	})
}

func TestGRPCClient_DeleteItem_NilItemClient(t *testing.T) {
	client := &GRPCClient{
		Item: nil,
	}

	var itemID [16]byte

	assert.Panics(t, func() {
		client.DeleteItem(context.Background(), "test-login", itemID)
	})
}

func TestGRPCClient_GetItems_NilClient(t *testing.T) {
	var client *GRPCClient = nil

	assert.Panics(t, func() {
		client.GetItems(context.Background(), "test-login", models.ItemTypeCARD)
	})
}

func TestGRPCClient_GetItems_NilItemClient(t *testing.T) {
	client := &GRPCClient{
		Item: nil,
	}

	assert.Panics(t, func() {
		client.GetItems(context.Background(), "test-login", models.ItemTypeCARD)
	})
}

func TestGRPCClient_GetTypesCounts_NilClient(t *testing.T) {
	var client *GRPCClient = nil

	assert.Panics(t, func() {
		client.GetTypesCounts(context.Background(), "test-login")
	})
}

func TestGRPCClient_GetTypesCounts_NilItemClient(t *testing.T) {
	client := &GRPCClient{
		Item: nil,
	}

	assert.Panics(t, func() {
		client.GetTypesCounts(context.Background(), "test-login")
	})
}
