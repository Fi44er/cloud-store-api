package user_usecase_contract

import (
	"context"

	user_entity "github.com/Fi44er/cloud-store-api/internal/modules/user/entity"
)

type IUserRepository interface {
	Create(ctx context.Context, user *user_entity.User) error
	Update(ctx context.Context, user *user_entity.User) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*user_entity.User, error)
	GetByEmail(ctx context.Context, email string) (*user_entity.User, error)
	GetByUsername(ctx context.Context, username string) (*user_entity.User, error)
	GetAll(ctx context.Context, offset, limit int) ([]user_entity.User, error)
}
