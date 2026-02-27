package file_http

import (
	"github.com/Fi44er/cloud-store-api/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func (h *NodeHandler) RegisterRoutes(router fiber.Router) {
	authMiddleware := middleware.NewAuthMiddleware()

	node := router.Group("/nodes", authMiddleware.RequireAuth)
	node.Post("/", h.CreateFolder)
}
