package activity_repository

import (
	file_entity "github.com/Fi44er/cloud-store-api/internal/modules/files/entity"
	file_model "github.com/Fi44er/cloud-store-api/internal/modules/files/infrastructure/repository/model"
)

type converter struct{}

func newConverter() *converter {
	return &converter{}
}

func (c *converter) toModel(log *file_entity.ActivityLog) *file_model.ActivityLog {
	var fileID *string
	if log.FileID != nil && *log.FileID != "" {
		fileID = log.FileID
	}

	return &file_model.ActivityLog{
		ID:        log.ID,
		UserID:    log.UserID,
		Action:    log.Action,
		FileID:    fileID,
		CreatedAt: log.CreatedAt,
	}
}

func (c *converter) toEntity(log *file_model.ActivityLog) *file_entity.ActivityLog {
	return &file_entity.ActivityLog{
		ID:        log.ID,
		UserID:    log.UserID,
		Action:    log.Action,
		FileID:    log.FileID,
		CreatedAt: log.CreatedAt,
	}
}
