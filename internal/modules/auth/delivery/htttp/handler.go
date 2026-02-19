package auth_http

import (
	auth_dto "github.com/Fi44er/cloud-store-api/internal/modules/auth/dto"
	auth_service "github.com/Fi44er/cloud-store-api/internal/modules/auth/service"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	logger *logger.Logger

	authService *auth_service.AuthService
}

func NewAuthHandler(logger *logger.Logger, authService *auth_service.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
	}
}

func (h *AuthHandler) InitRegistration(ctx *fiber.Ctx) error {
	flowID, err := h.authService.InitRegistration(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"flow_id": flowID,
	})
}

func (h *AuthHandler) Registration(ctx *fiber.Ctx) error {
	var dto auth_dto.RegistrationRequest
	if err := ctx.BodyParser(&dto); err != nil {
		return err
	}

	res, token, err := h.authService.Registration(ctx.Context(), &dto)
	if err != nil {
		if res != nil && len(res.Errors) > 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(res)
		}
		return err
	}

	if token != "" {
		ctx.Cookie(&fiber.Cookie{
			Name:     "ory_kratos_session",
			Value:    token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
		})
	}

	return ctx.Status(200).JSON(res)
}

func (h *AuthHandler) Verification(c *fiber.Ctx) error {
	var req auth_dto.VerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_body"})
	}

	cookie := c.Cookies("ory_kratos_session")

	result, err := h.authService.Verification(c.Context(), &req, cookie)
	if err != nil {
		return err
	}

	// Проверяем состояние потока
	// Если state == "passed_challenge", значит email подтвержден успешно
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"status":  result.VerificationFlow.State,
		"ui":      result.VerificationFlow.Ui, // здесь могут быть сообщения об успехе
	})
}

func (h *AuthHandler) InitLogin(ctx *fiber.Ctx) error {
	flowID, err := h.authService.InitLogin(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"flow_id": flowID,
	})
}

func (h *AuthHandler) Login(ctx *fiber.Ctx) error {
	var req auth_dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid_body"})
	}

	result, token, err := h.authService.Login(ctx.Context(), &req)
	if err != nil {
		return err
	}

	if token != "" {
		ctx.Cookie(&fiber.Cookie{
			Name:     "ory_kratos_session",
			Value:    token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
		})
	}

	return ctx.Status(200).JSON(result)
}

func (h *AuthHandler) WhoAmI(c *fiber.Ctx) error {
	session := c.Locals("session") // Предполагается, что middleware положило сюда сессию
	if session == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(session)
}
