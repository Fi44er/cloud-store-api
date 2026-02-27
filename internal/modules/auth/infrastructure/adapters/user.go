package auth_adapters

import (
	"context"

	auth_entity "github.com/Fi44er/cloud-store-api/internal/modules/auth/entity"
	user_entity "github.com/Fi44er/cloud-store-api/internal/modules/user/entity"
	user_usecase "github.com/Fi44er/cloud-store-api/internal/modules/user/usecase"
)

type UserUsecaseAdapter struct {
	userUseCase user_usecase.IUserUsecase
}

func NewUserUsecaseAdapter(userUsecase user_usecase.IUserUsecase) *UserUsecaseAdapter {
	return &UserUsecaseAdapter{
		userUseCase: userUsecase,
	}
}

func (a *UserUsecaseAdapter) Create(ctx context.Context, user *auth_entity.User) error {
	return a.userUseCase.Create(ctx, a.toUserEntity(user))
}

func (a *UserUsecaseAdapter) Update(ctx context.Context, user *auth_entity.User) error {
	return a.userUseCase.Update(ctx, a.toUserEntity(user))
}

func (a *UserUsecaseAdapter) Delete(ctx context.Context, id string) error {
	return a.userUseCase.Delete(ctx, id)
}

func (a *UserUsecaseAdapter) GetById(ctx context.Context, id string) (*auth_entity.User, error) {
	user, err := a.userUseCase.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return a.toAuthEntity(user), nil
}

func (a *UserUsecaseAdapter) GetByEmail(ctx context.Context, email string) (*auth_entity.User, error) {
	user, err := a.userUseCase.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return a.toAuthEntity(user), nil
}

func (a *UserUsecaseAdapter) GetByUsername(ctx context.Context, username string) (*auth_entity.User, error) {
	user, err := a.userUseCase.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return a.toAuthEntity(user), nil
}

func (a *UserUsecaseAdapter) GetAll(ctx context.Context, offset, limit int) ([]auth_entity.User, error) {
	users, err := a.userUseCase.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	var authUsers []auth_entity.User
	for _, user := range users {
		authUsers = append(authUsers, *a.toAuthEntity(&user))
	}
	return authUsers, nil
}

func (a *UserUsecaseAdapter) toAuthEntity(user *user_entity.User) *auth_entity.User {
	return &auth_entity.User{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
}

func (a *UserUsecaseAdapter) toUserEntity(user *auth_entity.User) *user_entity.User {
	return &user_entity.User{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
}
