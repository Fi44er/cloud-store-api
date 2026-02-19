package auth_module

import (
	auth_http "github.com/Fi44er/cloud-store-api/internal/modules/auth/delivery/htttp"
	auth_service "github.com/Fi44er/cloud-store-api/internal/modules/auth/service"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type AuthModule struct {
	logger *logger.Logger

	authService *auth_service.AuthService
	authHandler *auth_http.AuthHandler
}

func NewAuthModule(logger *logger.Logger) *AuthModule {
	return &AuthModule{
		logger: logger,
	}
}

func (m *AuthModule) Init() {
	m.authService = auth_service.NewAuthService(m.logger)
	m.authHandler = auth_http.NewAuthHandler(m.logger, m.authService)
}

func (m *AuthModule) InitDelivery(router fiber.Router) {
	m.authHandler.RegisterRoutes(router)
}
