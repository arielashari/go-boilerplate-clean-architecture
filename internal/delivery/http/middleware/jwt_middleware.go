package middleware

import (
	"strings"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/configs"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/delivery/http/response"
	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(cfg *configs.JWTConfig, authUseCase usecase.AuthUseCase) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token format", nil)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.Secret), nil
		})

		if err != nil || !token.Valid {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token", nil)
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["sub"].(string)
		tokenID := claims["jti"].(string)

		isValid, _ := authUseCase.ValidateSession(c.Context(), userID, tokenID)
		if !isValid {
			return response.Error(c, fiber.StatusUnauthorized, "Session revoked or expired", nil)
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
