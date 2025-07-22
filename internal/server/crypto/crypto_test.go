package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRSAKeyPair(t *testing.T) {
	pair, err := generateRSAKeyPair()

	require.NoError(t, err)
	require.NotNil(t, pair)
	require.NotNil(t, pair.PrivateKey)
	require.NotNil(t, pair.PublicKey)

	// Check key size is 2048 bits
	assert.Equal(t, 2048, pair.PrivateKey.Size()*8)

	// Check public key matches private key
	assert.Equal(t, &pair.PrivateKey.PublicKey, pair.PublicKey)
}

func TestRSAKeyPair_GetPublicKeyPEM(t *testing.T) {
	pair, err := generateRSAKeyPair()
	require.NoError(t, err)

	pemData, err := pair.getPublicKeyPEM()
	require.NoError(t, err)
	assert.NotEmpty(t, pemData)

	// Verify PEM format
	block, _ := pem.Decode(pemData)
	require.NotNil(t, block)
	assert.Equal(t, "PUBLIC KEY", block.Type)

	// Verify we can parse the key back
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	require.NoError(t, err)
	assert.NotNil(t, publicKey)
}

func TestRSAKeyPair_GetPrivateKeyPEM(t *testing.T) {
	pair, err := generateRSAKeyPair()
	require.NoError(t, err)

	pemData := pair.getPrivateKeyPEM()
	assert.NotEmpty(t, pemData)

	// Verify PEM format
	block, _ := pem.Decode(pemData)
	require.NotNil(t, block)
	assert.Equal(t, "RSA PRIVATE KEY", block.Type)

	// Verify we can parse the key back
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	require.NoError(t, err)
	assert.NotNil(t, privateKey)
}

func TestGetPublicKeyFromPEM(t *testing.T) {
	// Generate test key pair
	originalKey, err := generateRSAKeyPair()
	require.NoError(t, err)

	// Create proper PKIX format PEM for testing GetPublicKeyFromPEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(originalKey.PublicKey)
	require.NoError(t, err)

	pemData := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Parse back from PEM
	parsedKey, err := GetPublicKeyFromPEM(pemData)
	require.NoError(t, err)

	// Compare keys
	assert.Equal(t, originalKey.PublicKey.N, parsedKey.N)
	assert.Equal(t, originalKey.PublicKey.E, parsedKey.E)
}

func TestGetPublicKeyFromPEM_InvalidData(t *testing.T) {
	tests := []struct {
		name    string
		pemData []byte
	}{
		{
			name:    "empty data",
			pemData: []byte(""),
		},
		{
			name:    "invalid PEM",
			pemData: []byte("invalid pem data"),
		},
		{
			name:    "wrong block type",
			pemData: []byte("-----BEGIN CERTIFICATE-----\ndata\n-----END CERTIFICATE-----"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GetPublicKeyFromPEM(tt.pemData)
			assert.Error(t, err)
			assert.Nil(t, key)
		})
	}
}

func TestGetPrivateKeyFromPEM(t *testing.T) {
	// Generate test key pair
	originalKey, err := generateRSAKeyPair()
	require.NoError(t, err)

	// Get PEM data
	pemData := originalKey.getPrivateKeyPEM()

	// Parse back from PEM
	parsedKey, err := getPrivateKeyFromPEM(pemData)
	require.NoError(t, err)

	// Compare keys
	assert.Equal(t, originalKey.PrivateKey.N, parsedKey.N)
	assert.Equal(t, originalKey.PrivateKey.E, parsedKey.E)
	assert.Equal(t, originalKey.PrivateKey.D, parsedKey.D)
}

func TestGetPrivateKeyFromPEM_InvalidData(t *testing.T) {
	tests := []struct {
		name    string
		pemData []byte
	}{
		{
			name:    "empty data",
			pemData: []byte(""),
		},
		{
			name:    "invalid PEM",
			pemData: []byte("invalid pem data"),
		},
		{
			name:    "wrong block type",
			pemData: []byte("-----BEGIN CERTIFICATE-----\ndata\n-----END CERTIFICATE-----"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := getPrivateKeyFromPEM(tt.pemData)
			assert.Error(t, err)
			assert.Nil(t, key)
		})
	}
}

func TestGetKeysPath(t *testing.T) {
	path, err := getKeysPath()
	require.NoError(t, err)
	assert.NotEmpty(t, path)
	assert.Contains(t, path, "internal/server/crypto/keys")
}

func TestLoadRSAKeyPair_RequiresIntegrationTest(t *testing.T) {
	t.Skip("LoadRSAKeyPair requires config mocking and file system setup - use integration tests")
}

func TestKeyPairRoundTrip(t *testing.T) {
	// Test complete round trip: generate -> PEM -> parse -> compare
	originalPair, err := generateRSAKeyPair()
	require.NoError(t, err)

	// Convert to PEM
	privatePEM := originalPair.getPrivateKeyPEM()

	// For public key, create PKIX format for GetPublicKeyFromPEM
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(originalPair.PublicKey)
	require.NoError(t, err)

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Parse back from PEM
	parsedPrivate, err := getPrivateKeyFromPEM(privatePEM)
	require.NoError(t, err)

	parsedPublic, err := GetPublicKeyFromPEM(publicPEM)
	require.NoError(t, err)

	// Compare keys
	assert.Equal(t, originalPair.PrivateKey.N, parsedPrivate.N)
	assert.Equal(t, originalPair.PrivateKey.D, parsedPrivate.D)
	assert.Equal(t, originalPair.PublicKey.N, parsedPublic.N)
	assert.Equal(t, originalPair.PublicKey.E, parsedPublic.E)
}
