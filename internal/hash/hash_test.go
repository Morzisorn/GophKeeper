package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHash(t *testing.T) {
	password := []byte("testpassword")
	
	hash, err := GetHash(password)
	
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestVerifyHash(t *testing.T) {
	password := []byte("testpassword")
	
	hash, err := GetHash(password)
	require.NoError(t, err)
	
	// Correct password should return true
	result := VerifyHash(password, hash)
	assert.True(t, result)
	
	// Wrong password should return false
	wrongPassword := []byte("wrongpassword")
	result = VerifyHash(wrongPassword, hash)
	assert.False(t, result)
}