package user_module

import (
	user_repository "github.com/Fi44er/cloud-store-api/internal/modules/user/infrastructure/repository/user"
	user_usecase "github.com/Fi44er/cloud-store-api/internal/modules/user/usecase"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserModule struct {
	logger    *logger.Logger
	validator *validator.Validate
	db        *gorm.DB

	userRepo    user_repository.IUserRepository
	userUsecase user_usecase.IUserUsecase
}

func NewUserModule(logger *logger.Logger, validator *validator.Validate, db *gorm.DB) *UserModule {
	return &UserModule{
		logger:    logger,
		validator: validator,
		db:        db,
	}
}

func (m *UserModule) Init() {
	m.userRepo = user_repository.NewUserRepository(m.logger, m.db)
	m.userUsecase = user_usecase.NewUserUsecase(m.logger, m.userRepo)
}

func (m *UserModule) GetUserUseCase() user_usecase.IUserUsecase {
	return m.userUsecase
}
