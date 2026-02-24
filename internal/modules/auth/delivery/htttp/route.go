package auth_http

import (
	"github.com/Fi44er/cloud-store-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *AuthHandler) RegisterRoutes(router fiber.Router) {
	auth := router.Group("/auth")

	auth.Get("/registration/flow", h.InitRegistration)
	auth.Post("/registration", h.Registration)

	auth.Get("/login/flow", h.InitLogin)
	auth.Post("/login", h.Login)

	auth.Get("/verification/flow", h.InitVerification)
	auth.Post("/verification", h.Verification)
	auth.Post("/verification/resend", h.SendVerificationCode)

	authMiddleware := middleware.NewAuthMiddleware()

	account := router.Group("/account", authMiddleware.RequireAuth)
	account.Get("/me", h.WhoAmI)
	account.Post("/logout", h.Logout)

	session := account.Group("/session")
	session.Get("/", h.GetSessions)
	session.Delete("/:session_id", h.RevokeSession)
}
