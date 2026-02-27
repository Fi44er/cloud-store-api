package user_repository

import (
	user_entity "github.com/Fi44er/cloud-store-api/internal/modules/user/entity"
	user_model "github.com/Fi44er/cloud-store-api/internal/modules/user/infrastructure/repository/model"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
)

type Converter struct {
	logger *logger.Logger
}

func NewConverter(logger *logger.Logger) *Converter {
	return &Converter{
		logger: logger,
	}
}

func (c *Converter) ToModel(user *user_entity.User) *user_model.User {
	return &user_model.User{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		QuotaMax:   user.QuotaMax,
		AvatarPath: user.AvatarPath,
		BannerPath: user.BannerPath,
	}
}

func (c *Converter) ToEntity(user *user_model.User) *user_entity.User {
	return &user_entity.User{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		QuotaMax:   user.QuotaMax,
		AvatarPath: user.AvatarPath,
		BannerPath: user.BannerPath,
	}
}

func (c *Converter) ToEntities(users []user_model.User) []user_entity.User {
	var entities []user_entity.User
	for _, user := range users {
		entities = append(entities, *c.ToEntity(&user))
	}
	return entities
}
