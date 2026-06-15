package file_usecase

import (
	"context"
	"time"

	file_entity "github.com/Fi44er/cloud-store-api/internal/modules/files/entity"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
)

type ActivityUseCase struct {
	logger *logger.Logger
	repo   IActivityRepository
}

func NewActivityUseCase(logger *logger.Logger, repo IActivityRepository) *ActivityUseCase {
	return &ActivityUseCase{
		logger: logger,
		repo:   repo,
	}
}

func (u *ActivityUseCase) GetActivity(ctx context.Context, userID string) ([]file_entity.ActivityByDay, error) {
	from := time.Now().AddDate(-1, 0, 0)
	return u.repo.GetActivityByDays(ctx, userID, from)
}
