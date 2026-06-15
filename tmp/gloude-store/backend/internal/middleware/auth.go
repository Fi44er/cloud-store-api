// Пакет middleware содержит промежуточные обработчики запросов
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Claims — структура JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthRequired — middleware для проверки JWT токена
// Кладет user_id в c.Locals("user_id") для использования в хендлерах
// Документация Fiber Middleware: https://docs.gofiber.io/api/middleware
func AuthRequired(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем токен из заголовка Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			// Пробуем получить из cookie
			authHeader = "Bearer " + c.Cookies("token")
		}

		// Проверяем формат "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized: missing or invalid token",
			})
		}

		tokenStr := parts[1]

		// Парсим и валидируем JWT токен
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			// Проверяем алгоритм подписи
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized: invalid or expired token",
			})
		}

		// Кладем user_id в контекст запроса для использования в хендлерах
		// c.Locals(): https://docs.gofiber.io/api/ctx#locals
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}

// GetUserID извлекает user_id из контекста запроса
func GetUserID(c *fiber.Ctx) uint {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return 0
	}
	return userID
}
