package auth_module

import (
	auth_http "github.com/Fi44er/cloud-store-api/internal/modules/auth/delivery/htttp"
	auth_adapters "github.com/Fi44er/cloud-store-api/internal/modules/auth/infrastructure/adapters"
	auth_service "github.com/Fi44er/cloud-store-api/internal/modules/auth/service"
	user_usecase "github.com/Fi44er/cloud-store-api/internal/modules/user/usecase"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type AuthModule struct {
	logger *logger.Logger

	userUseCase user_usecase.IUserUsecase
	authService *auth_service.AuthService
	authHandler *auth_http.AuthHandler
}

func NewAuthModule(logger *logger.Logger) *AuthModule {
	return &AuthModule{
		logger: logger,
	}
}

func (m *AuthModule) Init() {
	adapter := auth_adapters.NewUserUsecaseAdapter(m.userUseCase)
	m.authService = auth_service.NewAuthService(m.logger, adapter)
	m.authHandler = auth_http.NewAuthHandler(m.logger, m.authService)
}

func (m *AuthModule) SetUserUseCase(userUseCase user_usecase.IUserUsecase) {
	m.userUseCase = userUseCase
}

func (m *AuthModule) InitDelivery(router fiber.Router) {
	m.authHandler.RegisterRoutes(router)
}
