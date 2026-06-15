package node_repository

import (
	"context"

	file_entity "github.com/Fi44er/cloud-store-api/internal/modules/files/entity"
	file_model "github.com/Fi44er/cloud-store-api/internal/modules/files/infrastructure/repository/model"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"gorm.io/gorm"
)

type NodeRepository struct {
	logger    *logger.Logger
	converter *converter
	db        *gorm.DB
}

func NewNodeRepository(logger *logger.Logger, db *gorm.DB) *NodeRepository {
	return &NodeRepository{
		logger:    logger,
		converter: newConverter(),
		db:        db,
	}
}

func (r *NodeRepository) Create(ctx context.Context, node *file_entity.Node) error {
	r.logger.Info("Creating node")
	model := r.converter.toModel(node)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Errorf("Failed to create node: %v", err)
		return err
	}
	node.ID = model.ID
	r.logger.Info("Node created successfully")
	return nil
}

func (r *NodeRepository) FindByID(ctx context.Context, id string) (*file_entity.Node, error) {
	r.logger.Infof("Finding node by ID: %s", id)
	var model file_model.Node
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		r.logger.Errorf("Failed to find node by ID: %v", err)
		return nil, err
	}
	node := r.converter.toEntity(&model)
	r.logger.Info("Node found successfully")
	return node, nil
}

func (r *NodeRepository) ExistsByName(ctx context.Context, name, userID string, parentID *string) (bool, error) {
	r.logger.Infof("Checking if node exists by name: %s", name)
	var count int64
	query := r.db.WithContext(ctx).Model(&file_model.Node{}).
		Where("user_id = ? AND name = ?", userID, name)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", parentID)
	}

	if err := query.Count(&count).Error; err != nil {
		r.logger.Errorf("Failed to check if node exists by name: %v", err)
		return false, err
	}

	r.logger.Info("Node exists check completed")
	return count > 0, nil
}
func (r *NodeRepository) Delete(ctx context.Context, id, userID string) error {
	r.logger.Infof("Deleting node: %s for user: %s", id, userID)
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&file_model.Node{}).Error
}

func (r *NodeRepository) GetUsedSpace(ctx context.Context, userID string) (int64, error) {
	var used int64
	err := r.db.WithContext(ctx).
		Model(&file_model.Node{}).
		Where("user_id = ? AND is_dir = false", userID).
		Select("COALESCE(SUM(size), 0)").
		Scan(&used).Error
	return used, err
}

func (r *NodeRepository) UpdateFavorite(ctx context.Context, id, userID string, isFavorite bool) error {
	return r.db.WithContext(ctx).
		Model(&file_model.Node{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_favorite", isFavorite).Error
}

type FileFilter struct {
	Extension string
	MinSize   int64
	MaxSize   int64
	Search    string
	Page      int
	Limit     int
}

func (r *NodeRepository) FindByUser(ctx context.Context, userID string, filter FileFilter) ([]file_entity.Node, int64, error) {
	var models []file_model.Node
	var total int64

	query := r.db.WithContext(ctx).Model(&file_model.Node{}).Where("user_id = ?", userID)

	if filter.Extension != "" {
		query = query.Where("extension = ?", filter.Extension)
	}

	if filter.MinSize > 0 {
		query = query.Where("size >= ?", filter.MinSize)
	}
	if filter.MaxSize > 0 {
		query = query.Where("size <= ?", filter.MaxSize)
	}

	if filter.Search != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]file_entity.Node, len(models))
	for i, m := range models {
		entities[i] = *r.converter.toEntity(&m)
	}

	return entities, total, nil
}
