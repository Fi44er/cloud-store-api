package file_http

import (
	"github.com/Fi44er/cloud-store-api/internal/middleware"
	node_repository "github.com/Fi44er/cloud-store-api/internal/modules/files/infrastructure/repository/node"
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

func (h *NodeHandler) RegisterRoutes(router fiber.Router) {
	authMiddleware := middleware.NewAuthMiddleware()
	files := router.Group("/files", authMiddleware.RequireAuth)

	files.Post("/folder", h.CreateFolder)
	files.Post("/upload", h.Upload)
	files.Get("/", h.ListFiles)
	files.Delete("/:id", h.DeleteNode)
	files.Get("/quota", h.GetQuota)
	files.Post("/favorite/:id", h.ToggleFavorite)
	files.Get("/download/:id", h.Download)
}

// Download godoc
// @Summary Download a file
// @Description Download a file from storage by ID
// @Tags files
// @Param id path string true "Node ID"
// @Success 200 {file} file
// @Failure 404 {object} map[string]string
// @Router /files/download/{id} [get]
func (h *NodeHandler) Download(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("identity_id").(string)

	node, err := h.useCase.GetNode(c.Context(), userID, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "file not found"})
	}

	if node.IsDir {
		return c.Status(400).JSON(fiber.Map{"error": "cannot download a directory"})
	}

	c.Set("Content-Disposition", `attachment; filename="`+node.Name+`"`)
	c.Set("Content-Type", node.MimeType)

	return c.SendFile(node.StoragePath, false)
}

// CreateFolder godoc
// @Summary Create a folder
// @Description Create a new folder in the storage
// @Tags files
// @Accept json
// @Produce json
// @Param request body createFolderRequest true "Folder creation request"
// @Success 201 {object} file_entity.Node
// @Failure 400 {object} map[string]string
// @Router /files/folder [post]
func (h *NodeHandler) CreateFolder(c *fiber.Ctx) error {
	var req createFolderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	userID := c.Locals("identity_id").(string)

	folder, err := h.useCase.CreateFolder(c.Context(), userID, req.Name, req.ParentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(folder)
}

// Upload godoc
// @Summary Upload a file
// @Description Upload a file to the storage
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param parent_id formData string false "Parent folder ID"
// @Success 201 {object} file_entity.Node
// @Failure 400 {object} map[string]string
// @Router /files/upload [post]
func (h *NodeHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "file is required"})
	}

	parentIDStr := c.FormValue("parent_id")
	var parentID *string
	if parentIDStr != "" {
		parentID = &parentIDStr
	}

	userID := c.Locals("identity_id").(string)

	node, err := h.useCase.Upload(c.Context(), userID, file, parentID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(node)
}

// ListFiles godoc
// @Summary List files
// ListFiles godoc
// @Summary List files and folders
// @Description Get a list of user files with filters
// @Tags files
// @Produce json
// @Param extension query string false "Filter by extension"
// @Param search query string false "Search by name"
// @Param min_size query int false "Minimum file size"
// @Param max_size query int false "Maximum file size"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /files [get]
func (h *NodeHandler) ListFiles(c *fiber.Ctx) error {
	userID := c.Locals("identity_id").(string)

	filter := node_repository.FileFilter{
		Extension: c.Query("extension"),
		Search:    c.Query("search"),
		MinSize:   int64(c.QueryInt("min_size", 0)),
		MaxSize:   int64(c.QueryInt("max_size", 0)),
		Page:      c.QueryInt("page", 1),
		Limit:     c.QueryInt("limit", 20),
	}

	files, total, err := h.useCase.ListFiles(c.Context(), userID, filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"files": files,
		"total": total,
	})
}

// DeleteNode godoc
// @Summary Delete a file or folder
// @Description Delete a file or folder from storage
// @Tags files
// @Param id path string true "Node ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Router /files/{id} [delete]
func (h *NodeHandler) DeleteNode(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("identity_id").(string)

	if err := h.useCase.DeleteNode(c.Context(), userID, id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

// GetQuota godoc
// @Summary Get storage quota
// @Description Get user's storage usage information
// @Tags files
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /files/quota [get]
func (h *NodeHandler) GetQuota(c *fiber.Ctx) error {
	userID := c.Locals("identity_id").(string)

	quota, err := h.useCase.GetQuota(c.Context(), userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(quota)
}

// ToggleFavorite godoc
// @Summary Toggle favorite status
// @Description Set or unset file as favorite
// @Tags files
// @Produce json
// @Param id path string true "Node ID"
// @Param is_favorite body bool true "Favorite status"
// @Success 200 "OK"
// @Failure 400 {object} map[string]string
// @Router /files/favorite/{id} [post]
func (h *NodeHandler) ToggleFavorite(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("identity_id").(string)

	var req struct {
		IsFavorite bool `json:"is_favorite"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := h.useCase.ToggleFavorite(c.Context(), userID, id, req.IsFavorite); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(200)
}
