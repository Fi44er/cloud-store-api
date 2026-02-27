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
