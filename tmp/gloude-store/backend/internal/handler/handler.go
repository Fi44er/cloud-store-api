// Пакет handler содержит обработчики HTTP запросов (контроллеры)
// Документация Fiber: https://docs.gofiber.io/
package handler

import (
	"errors"
	"strconv"

	"github.com/gloude/store/internal/middleware"
	"github.com/gloude/store/internal/models"
	"github.com/gloude/store/internal/service"
	"github.com/gofiber/fiber/v2"
)

// StorageHandler — обработчики эндпоинтов хранилища
type StorageHandler struct {
	storageSvc service.StorageService
}

// AuthHandler — обработчики эндпоинтов аутентификации
type AuthHandler struct {
	authSvc service.AuthService
}

// NewStorageHandler создает новый экземпляр обработчика хранилища
func NewStorageHandler(storageSvc service.StorageService) *StorageHandler {
	return &StorageHandler{storageSvc: storageSvc}
}

// NewAuthHandler создает новый экземпляр обработчика аутентификации
func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// ===== AUTH HANDLERS =====

// Register обрабатывает регистрацию нового пользователя
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Валидация входных данных
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "email, username and password are required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password must be at least 6 characters",
		})
	}

	user, err := h.authSvc.Register(c.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "user with this email already exists",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user registered successfully",
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// Login обрабатывает вход пользователя
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	user, token, err := h.authSvc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid email or password",
		})
	}

	// Устанавливаем JWT в httpOnly cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		SameSite: "lax",
		MaxAge:   86400, // 24 часа
	})

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

// Logout выходит из системы
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Очищаем cookie
	c.ClearCookie("token")
	return c.JSON(fiber.Map{"message": "logged out successfully"})
}

// Me возвращает данные текущего пользователя
// GET /api/v1/account/me
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	user, err := h.authSvc.GetUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(fiber.Map{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"quota_max":  user.QuotaMax,
		"created_at": user.CreatedAt,
	})
}

// ===== STORAGE HANDLERS =====

// Upload обрабатывает загрузку файла
// POST /api/v1/storage/upload
// Документация Fiber File Upload: https://docs.gofiber.io/api/ctx#formfile
func (h *StorageHandler) Upload(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Получаем файл из запроса
	// FormFile: https://docs.gofiber.io/api/ctx#formfile
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file is required in 'file' field",
		})
	}

	// Вызываем сервис для загрузки файла
	file, err := h.storageSvc.Upload(c.Context(), userID, fileHeader)
	if err != nil {
		if errors.Is(err, service.ErrQuotaExceeded) {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": "quota exceeded: not enough storage space",
			})
		}
		if errors.Is(err, service.ErrInvalidFileName) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid file name",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload file",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(file)
}

// ListFiles возвращает список файлов пользователя с фильтрацией
// GET /api/v1/storage/files
func (h *StorageHandler) ListFiles(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	// Парсим параметры фильтрации из query string
	filter := models.FileFilter{
		Extension: c.Query("extension"),
		Search:    c.Query("search"),
	}

	// Парсим числовые параметры
	if minSize := c.Query("min_size"); minSize != "" {
		filter.MinSize, _ = strconv.ParseInt(minSize, 10, 64)
	}
	if maxSize := c.Query("max_size"); maxSize != "" {
		filter.MaxSize, _ = strconv.ParseInt(maxSize, 10, 64)
	}
	if page := c.Query("page"); page != "" {
		filter.Page, _ = strconv.Atoi(page)
	}
	if limit := c.Query("limit"); limit != "" {
		filter.Limit, _ = strconv.Atoi(limit)
	}

	files, total, err := h.storageSvc.ListFiles(c.Context(), userID, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get files",
		})
	}

	return c.JSON(fiber.Map{
		"files": files,
		"total": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	})
}

// Download скачивает файл по ID (стриминг)
// GET /api/v1/storage/download/:file_id
// Документация Fiber Static: https://docs.gofiber.io/api/ctx#sendfile
func (h *StorageHandler) Download(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	fileID, err := strconv.ParseUint(c.Params("file_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid file id",
		})
	}

	// Получаем метаданные файла
	file, err := h.storageSvc.GetFile(c.Context(), userID, uint(fileID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "file not found",
		})
	}

	// Устанавливаем заголовки для скачивания
	c.Set("Content-Disposition", `attachment; filename="`+file.OriginalName+`"`)
	c.Set("Content-Type", file.MimeType)

	// Отправляем файл через стриминг
	// SendFile: https://docs.gofiber.io/api/ctx#sendfile
	return c.SendFile(file.Path, false)
}

// Delete удаляет файл
// DELETE /api/v1/storage/:file_id
func (h *StorageHandler) Delete(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	fileID, err := strconv.ParseUint(c.Params("file_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid file id",
		})
	}

	if err := h.storageSvc.DeleteFile(c.Context(), userID, uint(fileID)); err != nil {
		if errors.Is(err, service.ErrFileNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "file not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete file",
		})
	}

	return c.JSON(fiber.Map{"message": "file deleted successfully"})
}

// GetQuota возвращает информацию о квоте пользователя
// GET /api/v1/storage/quota
func (h *StorageHandler) GetQuota(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	quota, err := h.storageSvc.GetQuota(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get quota info",
		})
	}

	return c.JSON(quota)
}

// GetActivity возвращает данные активности для heatmap
// GET /api/v1/storage/activity
func (h *StorageHandler) GetActivity(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	activity, err := h.storageSvc.GetActivity(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get activity data",
		})
	}

	return c.JSON(fiber.Map{
		"activity": activity,
	})
}

// ToggleFavorite добавляет/убирает файл из избранного
// PATCH /api/v1/storage/:file_id/favorite
func (h *StorageHandler) ToggleFavorite(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	fileID, err := strconv.ParseUint(c.Params("file_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid file id",
		})
	}

	var req struct {
		IsFavorite bool `json:"is_favorite"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.storageSvc.ToggleFavorite(c.Context(), userID, uint(fileID), req.IsFavorite); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update favorite status",
		})
	}

	return c.JSON(fiber.Map{"message": "updated", "is_favorite": req.IsFavorite})
}
