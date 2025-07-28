package crypto

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt()

	require.NoError(t, err)
	assert.NotEmpty(t, salt)

	decoded, err := base64.StdEncoding.DecodeString(salt)
	require.NoError(t, err)

	assert.Len(t, decoded, 32)
}

func TestGenerateSalt_UniqueValues(t *testing.T) {
	salts := make([]string, 10)

	for i := 0; i < 10; i++ {
		salt, err := GenerateSalt()
		require.NoError(t, err)
		salts[i] = salt
	}

	for i := 0; i < len(salts); i++ {
		for j := i + 1; j < len(salts); j++ {
			assert.NotEqual(t, salts[i], salts[j], "Salt %d and %d should be different", i, j)
		}
	}
}

func TestGenerateSalt_Length(t *testing.T) {
	salt, err := GenerateSalt()

	require.NoError(t, err)

	expectedLength := base64.StdEncoding.EncodedLen(32)
	assert.Len(t, salt, expectedLength)
}

func TestGenerateSalt_ValidBase64(t *testing.T) {
	salt, err := GenerateSalt()

	require.NoError(t, err)

	_, err = base64.StdEncoding.DecodeString(salt)
	assert.NoError(t, err, "Generated salt should be valid base64")
}

func TestGenerateSalt_MultipleCallsSucceed(t *testing.T) {
	for i := 0; i < 100; i++ {
		salt, err := GenerateSalt()
		assert.NoError(t, err)
		assert.NotEmpty(t, salt)
	}
}

func TestGenerateSalt_RandomnessCheck(t *testing.T) {
	salts := make(map[string]bool)
	const iterations = 1000

	for i := 0; i < iterations; i++ {
		salt, err := GenerateSalt()
		require.NoError(t, err)

		assert.False(t, salts[salt], "Duplicate salt generated at iteration %d", i)
		salts[salt] = true
	}

	assert.Len(t, salts, iterations)
}
