package auth_http

import (
	"github.com/Fi44er/cloud-store-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Get("/registration/flow", h.InitRegistration)
	auth.Post("/registration", h.Registration)
	auth.Post("/verification", h.Verification)

	auth.Get("/login/flow", h.InitLogin)
	auth.Post("/login", h.Login)

	authMiddleware := middleware.NewAuthMiddleware()

	protected := router.Group("/auth", authMiddleware.RequireAuth)
	protected.Get("/whoami", h.WhoAmI)
}
