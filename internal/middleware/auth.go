package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	kratos "github.com/ory/kratos-client-go"
)

type AuthMiddleware struct {
	kratosClient *kratos.APIClient
}

func NewAuthMiddleware() *AuthMiddleware {
	configuration := kratos.NewConfiguration()
	configuration.Servers = []kratos.ServerConfiguration{
		{
			URL: "http://localhost:4433",
		},
	}

	return &AuthMiddleware{
		kratosClient: kratos.NewAPIClient(configuration),
	}
}

func (m *AuthMiddleware) extractToken(c *fiber.Ctx) string {
	// 1. Пытаемся взять из куки
	token := c.Cookies("ory_kratos_session")
	if token != "" {
		return token
	}

	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return ""
}

func (m *AuthMiddleware) validateSession(c *fiber.Ctx, token string) (*kratos.Session, error) {
	session, resp, err := m.kratosClient.FrontendAPI.ToSession(c.UserContext()).
		XSessionToken(token).
		Execute()

	if err != nil || resp == nil || resp.StatusCode != 200 {
		return nil, err
	}

	if session == nil || session.Active == nil || !*session.Active {
		return nil, nil
	}

	return session, nil
}

func (m *AuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	token := m.extractToken(c)

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "unauthorized",
			"message": "Authentication required",
		})
	}

	session, err := m.validateSession(c, token)
	if err != nil || session == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "unauthorized",
			"message": "Invalid or expired session",
		})
	}

	c.Locals("session", session)
	c.Locals("identity_id", session.Identity.Id)

	return c.Next()
}

func (m *AuthMiddleware) OptionalAuth(c *fiber.Ctx) error {
	token := m.extractToken(c)
	if token == "" {
		return c.Next()
	}

	session, err := m.validateSession(c, token)
	if err == nil && session != nil {
		c.Locals("session", session)
		c.Locals("identity_id", session.Identity.Id)
	}

	return c.Next()
}
