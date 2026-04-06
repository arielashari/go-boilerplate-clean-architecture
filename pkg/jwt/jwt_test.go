package jwt_test

import (
	"testing"
	"time"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	_jwt "github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/jwt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func testConfig() *configs.JWTConfig {
	return &configs.JWTConfig{
		Secret:               "test-secret-key",
		AccessExpireMinutes:  15,
		RefreshExpireMinutes: 10080,
	}
}

func TestGenerateTokenPair(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		td, err := _jwt.GenerateTokenPair("user-1", testConfig())

		assert.NoError(t, err)
		assert.NotEmpty(t, td.AccessToken)
		assert.NotEmpty(t, td.RefreshToken)
		assert.NotEmpty(t, td.TokenID)
		assert.NotEqual(t, td.AccessToken, td.RefreshToken)
	})

	t.Run("access token has correct claims", func(t *testing.T) {
		td, err := _jwt.GenerateTokenPair("user-1", testConfig())
		assert.NoError(t, err)

		token, err := jwt.Parse(td.AccessToken, func(t *jwt.Token) (interface{}, error) {
			return []byte("test-secret-key"), nil
		})

		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims := token.Claims.(jwt.MapClaims)
		assert.Equal(t, "user-1", claims["sub"])
		assert.Equal(t, "access", claims["type"])
		assert.Equal(t, td.TokenID, claims["jti"])
	})

	t.Run("refresh token has correct claims", func(t *testing.T) {
		td, err := _jwt.GenerateTokenPair("user-1", testConfig())
		assert.NoError(t, err)

		token, err := jwt.Parse(td.RefreshToken, func(t *jwt.Token) (interface{}, error) {
			return []byte("test-secret-key"), nil
		})

		assert.NoError(t, err)
		assert.True(t, token.Valid)

		claims := token.Claims.(jwt.MapClaims)
		assert.Equal(t, "user-1", claims["sub"])
		assert.Equal(t, "refresh", claims["type"])
		assert.Equal(t, td.TokenID, claims["jti"])
	})

	t.Run("access token expires before refresh token", func(t *testing.T) {
		td, err := _jwt.GenerateTokenPair("user-1", testConfig())
		assert.NoError(t, err)
		assert.Less(t, td.AtExpires, td.RtExpires)
	})

	t.Run("different users get different token IDs", func(t *testing.T) {
		td1, _ := _jwt.GenerateTokenPair("user-1", testConfig())
		td2, _ := _jwt.GenerateTokenPair("user-2", testConfig())
		assert.NotEqual(t, td1.TokenID, td2.TokenID)
	})

	t.Run("token invalid with wrong secret", func(t *testing.T) {
		td, err := _jwt.GenerateTokenPair("user-1", testConfig())
		assert.NoError(t, err)

		_, err = jwt.Parse(td.AccessToken, func(t *jwt.Token) (interface{}, error) {
			return []byte("wrong-secret"), nil
		})
		assert.Error(t, err)
	})

	t.Run("expires at is in the future", func(t *testing.T) {
		td, err := _jwt.GenerateTokenPair("user-1", testConfig())
		assert.NoError(t, err)
		assert.Greater(t, td.AtExpires, time.Now().Unix())
		assert.Greater(t, td.RtExpires, time.Now().Unix())
	})
}
