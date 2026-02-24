package auth_utils

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetClientIP(c *fiber.Ctx) string {
	if debugIP := os.Getenv("APP_DEBUG_IP"); debugIP != "" {
		return debugIP
	}
	return c.IP()
}
