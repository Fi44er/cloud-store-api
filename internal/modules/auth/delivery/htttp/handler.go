package auth_http

import (
	auth_dto "github.com/Fi44er/cloud-store-api/internal/modules/auth/dto"
	auth_constant "github.com/Fi44er/cloud-store-api/internal/modules/auth/pkg/constant"
	auth_utils "github.com/Fi44er/cloud-store-api/internal/modules/auth/pkg/util"
	auth_service "github.com/Fi44er/cloud-store-api/internal/modules/auth/service"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	kratos "github.com/ory/kratos-client-go"
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

// =================== Registration ===================
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
			Name:     auth_constant.CratosSessionKey,
			Value:    token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
		})
	}

	return ctx.Status(200).JSON(res)
}

// =================== Verification ===================
func (h *AuthHandler) InitVerification(ctx *fiber.Ctx) error {
	flowID, err := h.authService.InitVerification(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"flow_id": flowID,
	})
}

func (h *AuthHandler) SendVerificationCode(ctx *fiber.Ctx) error {
	var req auth_dto.VerificationSendCodeRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid_body"})
	}

	err := h.authService.SendVerificationCode(ctx.Context(), req.FlowID, req.Email)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"success": true,
	})
}

func (h *AuthHandler) Verification(c *fiber.Ctx) error {
	var req auth_dto.VerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid_body"})
	}

	cookie := c.Cookies(auth_constant.CratosSessionKey)

	result, err := h.authService.Verification(c.Context(), &req, cookie)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{
		"status": result.Status,
	})
}

// =================== Login ===================
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
	ip := auth_utils.GetClientIP(ctx)
	userAgent := ctx.Get("User-Agent")
	var req auth_dto.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid_body"})
	}

	result, token, err := h.authService.Login(ctx.Context(), &req, userAgent, ip)
	if err != nil {
		return err
	}

	if token != "" {
		ctx.Cookie(&fiber.Cookie{
			Name:     auth_constant.CratosSessionKey,
			Value:    token,
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
			Path:     "/",
		})
	}

	return ctx.Status(200).JSON(result)
}

// =================== Logout ===================
func (h *AuthHandler) Logout(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies(auth_constant.CratosSessionKey)
	res, err := h.authService.Logout(ctx.Context(), cookie)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     auth_constant.CratosSessionKey,
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	return ctx.Status(200).JSON(res)
}

// =================== Sessions ===================
func (h *AuthHandler) WhoAmI(c *fiber.Ctx) error {
	session := c.Locals(auth_constant.SessionCtxKey).(*kratos.Session)
	if session == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(session)
}

func (h *AuthHandler) RevokeSession(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies(auth_constant.CratosSessionKey)
	sessionID := ctx.Params("session_id")
	err := h.authService.RevokeSession(ctx.Context(), cookie, sessionID)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{"success": true})
}

func (h *AuthHandler) GetSessions(ctx *fiber.Ctx) error {
	identityID := ctx.Locals(auth_constant.IdentityIdCtxKey).(string)
	session := ctx.Locals(auth_constant.SessionCtxKey).(*kratos.Session)
	sessions, err := h.authService.GetSessions(ctx.Context(), identityID, session.Id)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(sessions)
}
