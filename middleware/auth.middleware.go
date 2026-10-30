package middleware

import (
	"strings"

	"github.com/RaucousDave/streakly/utils"
	"github.com/gofiber/fiber/v3"
)

func RequireAuth(c fiber.Ctx) error {
	header := c.Get("Authorization")

	parts := strings.SplitN(header, " ", 2)

	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return utils.Unauthorized(c, "missing or malformed token", nil)
	}

	claims, err := utils.VerifyToken(parts[1])
	if err != nil || claims.TokenType != utils.TypeAccess {
		return utils.Unauthorized(c, "invalid or expired token", nil)
	}

	c.Locals("userId", claims.UserId)
	return c.Next()
}
