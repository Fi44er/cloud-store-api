package file_http

import (
	"fmt"

	file_usecase "github.com/Fi44er/cloud-store-api/internal/modules/files/usecase"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type NodeHandler struct {
	logger  *logger.Logger
	useCase *file_usecase.NodeUseCase
}

func NewNodeHandler(logger *logger.Logger, useCase *file_usecase.NodeUseCase) *NodeHandler {
	return &NodeHandler{
		logger:  logger,
		useCase: useCase,
	}
}

type createFolderRequest struct {
	Name     string  `json:"name" validate:"required"`
	ParentID *string `json:"parent_id"` // Берем строку, чтобы потом валидировать UUID
}

func (h *NodeHandler) CreateFolder(c *fiber.Ctx) error {
	fmt.Println("bam")
	var req createFolderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name is required"})
	}

	// Получаем UserID из контекста (предположим, его положил туда Auth Middleware)
	userID := c.Locals("identity_id").(string)

	// Вызов бизнес-логики
	folder, err := h.useCase.CreateFolder(c.Context(), userID, req.Name, req.ParentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(folder)
}
