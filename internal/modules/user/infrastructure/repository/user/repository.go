package user_repository

import (
	"context"

	user_entity "github.com/Fi44er/cloud-store-api/internal/modules/user/entity"
	user_model "github.com/Fi44er/cloud-store-api/internal/modules/user/infrastructure/repository/model"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(ctx context.Context, user *user_entity.User) error
	Update(ctx context.Context, user *user_entity.User) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*user_entity.User, error)
	GetByEmail(ctx context.Context, email string) (*user_entity.User, error)
	GetByUsername(ctx context.Context, username string) (*user_entity.User, error)
	GetAll(ctx context.Context, offset, limit int) ([]*user_entity.User, error)
}

type UserRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewUserRepository(logger *logger.Logger, db *gorm.DB) IUserRepository {
	return &UserRepository{
		logger:    logger,
		db:        db,
		converter: NewConverter(logger),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *user_entity.User) error {
	r.logger.Infof("Creating user %s", user.Email)
	userModel := r.converter.ToModel(user)
	err := r.db.WithContext(ctx).Create(&userModel).Error
	if err != nil {
		r.logger.Errorf("Failed to create user %s: %v", user.Email, err)
		return err
	}
	user.ID = userModel.ID
	r.logger.Infof("User %s created successfully", user.Email)
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *user_entity.User) error {
	r.logger.Infof("Updating user %s", user.Email)
	userModel := r.converter.ToModel(user)
	err := r.db.WithContext(ctx).Model(&userModel).Updates(userModel).Error
	if err != nil {
		r.logger.Errorf("Failed to update user %s: %v", user.Email, err)
		return err
	}
	r.logger.Infof("User %s updated successfully", user.Email)
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, userID string) error {
	r.logger.Infof("Deleting user with ID %s", userID)
	err := r.db.WithContext(ctx).Delete(&user_entity.User{}, userID).Error
	if err != nil {
		r.logger.Errorf("Failed to delete user with ID %s: %v", userID, err)
		return err
	}
	r.logger.Infof("User with ID %s deleted successfully", userID)
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*user_entity.User, error) {
	r.logger.Infof("Getting user with ID %s", userID)
	var userModel user_model.User
	err := r.db.WithContext(ctx).First(&userModel, userID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("User not found: %s", userID)
			return nil, nil
		}
		r.logger.Errorf("Failed to get user with ID %s: %v", userID, err)
		return nil, err
	}
	user := r.converter.ToEntity(&userModel)
	r.logger.Infof("User with ID %s retrieved successfully", userID)
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user_entity.User, error) {
	r.logger.Infof("Getting user with email %s", email)
	var userModel user_model.User
	err := r.db.WithContext(ctx).First(&userModel, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("User not found: %s", email)
			return nil, nil
		}
		r.logger.Errorf("Failed to get user with email %s: %v", email, err)
		return nil, err
	}
	user := r.converter.ToEntity(&userModel)
	r.logger.Infof("User with email %s retrieved successfully", email)
	return user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*user_entity.User, error) {
	r.logger.Infof("Getting user with username %s", username)
	var userModel user_model.User
	err := r.db.WithContext(ctx).First(&userModel, "username = ?", username).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("User not found: %s", username)
			return nil, nil
		}
		r.logger.Errorf("Failed to get user with username %s: %v", username, err)
		return nil, err
	}
	user := r.converter.ToEntity(&userModel)
	r.logger.Infof("User with username %s retrieved successfully", username)
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, offset, limit int) ([]*user_entity.User, error) {
	r.logger.Infof("Getting all users with offset %d and limit %d", offset, limit)
	var userModels []user_model.User
	if limit == 0 {
		limit = -1
	}
	if offset == 0 {
		offset = -1
	}
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&userModels).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("No users found")
			return nil, nil
		}
		r.logger.Errorf("Failed to get all users: %v", err)
		return nil, err
	}
	users := r.converter.ToEntities(userModels)
	r.logger.Infof("All users retrieved successfully")
	return users, nil
}
