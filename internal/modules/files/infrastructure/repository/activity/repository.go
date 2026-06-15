package activity_repository

import (
	"context"
	"time"

	file_entity "github.com/Fi44er/cloud-store-api/internal/modules/files/entity"
	file_model "github.com/Fi44er/cloud-store-api/internal/modules/files/infrastructure/repository/model"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	logger    *logger.Logger
	converter *converter
	db        *gorm.DB
}

func NewActivityRepository(logger *logger.Logger, db *gorm.DB) *ActivityRepository {
	return &ActivityRepository{
		logger:    logger,
		converter: newConverter(),
		db:        db,
	}
}

func (r *ActivityRepository) Create(ctx context.Context, log *file_entity.ActivityLog) error {
	r.logger.Info("Creating activity log")
	model := r.converter.toModel(log)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Errorf("Failed to create activity log: %v", err)
		return err
	}
	return nil
}

func (r *ActivityRepository) GetActivityByDays(ctx context.Context, userID string, from time.Time) ([]file_entity.ActivityByDay, error) {
	var result []file_entity.ActivityByDay

	err := r.db.WithContext(ctx).
		Model(&file_model.ActivityLog{}).
		Select("TO_CHAR(DATE_TRUNC('day', created_at), 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("user_id = ? AND action = 'upload' AND created_at >= ?", userID, from).
		Group("DATE_TRUNC('day', created_at)").
		Order("date ASC").
		Scan(&result).Error

	return result, err
}
