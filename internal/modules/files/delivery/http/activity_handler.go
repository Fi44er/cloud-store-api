package file_http

import (
	"github.com/Fi44er/cloud-store-api/internal/middleware"
	file_usecase "github.com/Fi44er/cloud-store-api/internal/modules/files/usecase"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type ActivityHandler struct {
	logger  *logger.Logger
	useCase *file_usecase.ActivityUseCase
}

func NewActivityHandler(logger *logger.Logger, useCase *file_usecase.ActivityUseCase) *ActivityHandler {
	return &ActivityHandler{
		logger:  logger,
		useCase: useCase,
	}
}

func (h *ActivityHandler) RegisterRoutes(router fiber.Router) {
	authMiddleware := middleware.NewAuthMiddleware()
	router.Get("/files/activity", authMiddleware.RequireAuth, h.GetActivity)
}

// GetActivity godoc
// @Summary Get activity heatmap
// @Description Get user's file activity for the last 365 days
// @Tags files
// @Produce json
// @Success 200 {array} file_entity.ActivityByDay
// @Router /files/activity [get]
func (h *ActivityHandler) GetActivity(c *fiber.Ctx) error {
	userID := c.Locals("identity_id").(string)

	activity, err := h.useCase.GetActivity(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(activity)
}
