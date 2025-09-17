package middleware

import (
	"strings"

	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct {
	jwtSecret string
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{jwtSecret: jwtSecret}
}

func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Authorization header required", "")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization format", "")
		}

		userID, err := auth.ValidateToken(tokenString, m.jwtSecret)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid token", err.Error())
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
