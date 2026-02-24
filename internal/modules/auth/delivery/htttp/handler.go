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

// InitRegistration initiates a new registration flow
// @Summary Initialize registration flow
// @Description Starts a new registration flow and returns a flow ID
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "flow_id"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/registration/flow [get]
func (h *AuthHandler) InitRegistration(ctx *fiber.Ctx) error {
	flowID, err := h.authService.InitRegistration(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"flow_id": flowID,
	})
}

// Registration registers a new user
// @Summary Register a new user
// @Description Registers a new user with email, username, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth_dto.RegistrationRequest true "Registration details"
// @Success 200 {object} auth_dto.RegisterResponse "Registration successful"
// @Failure 400 {object} auth_dto.RegisterResponse "Validation error"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/registration [post]
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

// InitVerification initiates a new verification flow
// @Summary Initialize verification flow
// @Description Starts a new verification flow and returns a flow ID
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "flow_id"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/verification/flow [get]
func (h *AuthHandler) InitVerification(ctx *fiber.Ctx) error {
	flowID, err := h.authService.InitVerification(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"flow_id": flowID,
	})
}

// SendVerificationCode sends a verification code to the user's email
// @Summary Send verification code
// @Description Sends a verification code to the specified email address
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth_dto.VerificationSendCodeRequest true "Email and flow ID"
// @Success 200 {object} map[string]bool "success"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/verification/resend [post]
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

// Verification verifies a user with a code
// @Summary Verify user
// @Description Verifies a user using the provided verification code
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth_dto.VerificationRequest true "Verification code and flow ID"
// @Success 200 {object} auth_dto.VerificationResponse "Verification status"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/verification [post]
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

// InitLogin initiates a new login flow
// @Summary Initialize login flow
// @Description Starts a new login flow and returns a flow ID
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "flow_id"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login/flow [get]
func (h *AuthHandler) InitLogin(ctx *fiber.Ctx) error {
	flowID, err := h.authService.InitLogin(ctx.Context())
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"flow_id": flowID,
	})
}

// Login authenticates a user
// @Summary User login
// @Description Authenticates a user and creates a session
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth_dto.LoginRequest true "Login credentials"
// @Success 200 {object} auth_dto.LoginResponse "Login successful"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} auth_dto.LoginResponse "Authentication failed"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
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

	h.logger.Debugf("token: %s %s", token, auth_constant.CratosSessionKey)
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

// Logout terminates the current session
// @Summary User logout
// @Description Logs out the current user by invalidating their session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} auth_dto.LogoutResponse "Logout successful"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /account/logout [post]
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

// WhoAmI returns the current session information
// @Summary Get current user session
// @Description Returns the session details of the authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} kratos.Session "Session details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /account/me [get]
func (h *AuthHandler) WhoAmI(c *fiber.Ctx) error {
	session := c.Locals(auth_constant.SessionCtxKey).(*kratos.Session)
	if session == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	return c.JSON(session)
}

// RevokeSession terminates a specific session
// @Summary Revoke user session
// @Description Revokes a specific session by ID for the authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param session_id path string true "Session ID to revoke"
// @Success 200 {object} map[string]bool "success"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /account/session/{session_id} [delete]
func (h *AuthHandler) RevokeSession(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies(auth_constant.CratosSessionKey)
	sessionID := ctx.Params("session_id")
	err := h.authService.RevokeSession(ctx.Context(), cookie, sessionID)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{"success": true})
}

// GetSessions returns all active sessions for the user
// @Summary Get user sessions
// @Description Returns all active sessions for the authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} auth_dto.SessionResponse "List of sessions"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /account/session [get]
func (h *AuthHandler) GetSessions(ctx *fiber.Ctx) error {
	identityID := ctx.Locals(auth_constant.IdentityIdCtxKey).(string)
	session := ctx.Locals(auth_constant.SessionCtxKey).(*kratos.Session)
	sessions, err := h.authService.GetSessions(ctx.Context(), identityID, session.Id)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(sessions)
}
