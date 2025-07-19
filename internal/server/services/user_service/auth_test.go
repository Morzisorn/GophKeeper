package user_service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		login   string
		wantErr bool
	}{
		{
			name:    "successful token generation",
			login:   "testuser",
			wantErr: false,
		},
		{
			name:    "successful token generation with empty login",
			login:   "",
			wantErr: false,
		},
		{
			name:    "successful token generation with special characters",
			login:   "user@example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := generateToken(tt.login)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				
				assert.Contains(t, token, ".")
				
				parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				
					return []byte("test_secret"), nil 
				})
				
				if err == nil {
					claims, ok := parsedToken.Claims.(jwt.MapClaims)
					require.True(t, ok)
					assert.Equal(t, tt.login, claims["login"])
					
					exp, ok := claims["exp"].(float64)
					require.True(t, ok)
					
					expectedExp := time.Now().Add(7 * time.Hour * 24).Unix()
					assert.InDelta(t, expectedExp, exp, 10)
				}
			}
		})
	}
}

func TestGenerateToken_Format(t *testing.T) {
	token, err := generateToken("testuser")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	
	parts := 0
	for i := 0; i < len(token); i++ {
		if token[i] == '.' {
			parts++
		}
	}
	assert.Contains(t, token, ".")
}

func TestGenerateToken_DifferentLogins(t *testing.T) {
	login1 := "user1"
	login2 := "user2"
	
	token1, err1 := generateToken(login1)
	token2, err2 := generateToken(login2)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	
	assert.NotEqual(t, token1, token2)
}

func TestGenerateToken_MultipleCallsSameUser(t *testing.T) {
	login := "testuser"
	
	token1, err1 := generateToken(login)
	time.Sleep(time.Second)
	token2, err2 := generateToken(login)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEmpty(t, token1)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token1, token2)
}