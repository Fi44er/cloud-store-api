package user_usecase

import (
	"context"

	user_entity "github.com/Fi44er/cloud-store-api/internal/modules/user/entity"
	user_usecase_contract "github.com/Fi44er/cloud-store-api/internal/modules/user/usecase/contract"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
)

type IUserUsecase interface {
	Create(ctx context.Context, user *user_entity.User) error
	GetByID(ctx context.Context, id string) (*user_entity.User, error)
	GetByEmail(ctx context.Context, email string) (*user_entity.User, error)
	GetByUsername(ctx context.Context, username string) (*user_entity.User, error)
	GetAll(ctx context.Context, offset, limit int) ([]*user_entity.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user *user_entity.User) error
}

type UserUsecase struct {
	logger *logger.Logger
	repo   user_usecase_contract.IUserRepository
}

func NewUserUsecase(logger *logger.Logger, repo user_usecase_contract.IUserRepository) IUserUsecase {
	return &UserUsecase{
		logger: logger,
		repo:   repo,
	}
}

func (u *UserUsecase) Create(ctx context.Context, user *user_entity.User) error {
	u.logger.Infof("Creating user: %v", user)
	if err := user.Validate(); err != nil {
		u.logger.Errorf("Failed to validate user: %v", err)
		return err
	}

	return u.repo.Create(ctx, user)
}

func (u *UserUsecase) GetByID(ctx context.Context, id string) (*user_entity.User, error) {
	u.logger.Infof("Getting user by ID: %s", id)
	return u.repo.GetByID(ctx, id)
}

func (u *UserUsecase) Update(ctx context.Context, user *user_entity.User) error {
	u.logger.Infof("Updating user: %v", user)
	if err := user.Validate(); err != nil {
		u.logger.Errorf("Failed to validate user: %v", err)
		return err
	}

	return u.repo.Update(ctx, user)
}

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	u.logger.Infof("Deleting user by ID: %s", id)
	return u.repo.Delete(ctx, id)
}

func (u *UserUsecase) GetByUsername(ctx context.Context, username string) (*user_entity.User, error) {
	u.logger.Infof("Getting user by username: %s", username)
	return u.repo.GetByUsername(ctx, username)
}

func (u *UserUsecase) GetByEmail(ctx context.Context, email string) (*user_entity.User, error) {
	u.logger.Infof("Getting user by email: %s", email)
	return u.repo.GetByEmail(ctx, email)
}

func (u *UserUsecase) GetAll(ctx context.Context, offset, limit int) ([]*user_entity.User, error) {
	u.logger.Infof("Getting all users with offset: %d and limit: %d", offset, limit)
	return u.repo.GetAll(ctx, offset, limit)
}
