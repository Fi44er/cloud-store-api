package file_usecase

import (
	"context"
	"errors"

	file_entity "github.com/Fi44er/cloud-store-api/internal/modules/files/entity"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
)

type INodeRepository interface {
	ExistsByName(ctx context.Context, name, userID string, parentID *string) (bool, error)
	FindByID(ctx context.Context, id string) (*file_entity.Node, error)
	Create(ctx context.Context, node *file_entity.Node) error
}

type NodeUseCase struct {
	logger *logger.Logger
	repo   INodeRepository
}

func NewNodeUseCase(logger *logger.Logger, repo INodeRepository) *NodeUseCase {
	return &NodeUseCase{
		logger: logger,
		repo:   repo,
	}
}

func (u *NodeUseCase) CreateFolder(ctx context.Context, userID, name string, parentID *string) (*file_entity.Node, error) {
	if parentID != nil {
		parent, err := u.repo.FindByID(ctx, *parentID)
		if err != nil {
			return nil, errors.New("parent folder not found")
		}
		if !parent.IsDir {
			return nil, errors.New("parent is not a directory")
		}
		if parent.UserID != userID {
			return nil, errors.New("access denied")
		}
	}

	exists, err := u.repo.ExistsByName(ctx, name, userID, parentID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("item with this name already exists")
	}

	folder := &file_entity.Node{
		Name:     name,
		IsDir:    true,
		UserID:   userID,
		ParentID: parentID,
	}

	if err := u.repo.Create(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}
