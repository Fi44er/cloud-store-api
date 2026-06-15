package file_usecase

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	file_entity "github.com/Fi44er/cloud-store-api/internal/modules/files/entity"
	node_repository "github.com/Fi44er/cloud-store-api/internal/modules/files/infrastructure/repository/node"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/google/uuid"
)

type INodeRepository interface {
	ExistsByName(ctx context.Context, name, userID string, parentID *string) (bool, error)
	FindByID(ctx context.Context, id string) (*file_entity.Node, error)
	Create(ctx context.Context, node *file_entity.Node) error
	Delete(ctx context.Context, id, userID string) error
	GetUsedSpace(ctx context.Context, userID string) (int64, error)
	UpdateFavorite(ctx context.Context, id, userID string, isFavorite bool) error
	FindByUser(ctx context.Context, userID string, filter node_repository.FileFilter) ([]file_entity.Node, int64, error)
}

type IActivityRepository interface {
	Create(ctx context.Context, log *file_entity.ActivityLog) error
	GetActivityByDays(ctx context.Context, userID string, from time.Time) ([]file_entity.ActivityByDay, error)
}

type NodeUseCase struct {
	logger       *logger.Logger
	repo         INodeRepository
	activityRepo IActivityRepository
	uploadDir    string
	maxQuota     int64
}

func NewNodeUseCase(logger *logger.Logger, repo INodeRepository, activityRepo IActivityRepository, uploadDir string, maxQuota int64) *NodeUseCase {
	return &NodeUseCase{
		logger:       logger,
		repo:         repo,
		activityRepo: activityRepo,
		uploadDir:    uploadDir,
		maxQuota:     maxQuota,
	}
}

func (u *NodeUseCase) CreateFolder(ctx context.Context, userID, name string, parentID *string) (*file_entity.Node, error) {
	node := &file_entity.Node{
		Name:     name,
		IsDir:    true,
		UserID:   userID,
		ParentID: parentID,
	}

	if err := u.repo.Create(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

func (u *NodeUseCase) Upload(ctx context.Context, userID string, fileHeader *multipart.FileHeader, parentID *string) (*file_entity.Node, error) {
	// 1. Check quota
	used, err := u.repo.GetUsedSpace(ctx, userID)
	if err != nil {
		return nil, err
	}
	if used+fileHeader.Size > u.maxQuota {
		return nil, errors.New("quota exceeded")
	}

	// 2. Open file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// 3. Prepare storage
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	storageName := uuid.New().String() + ext
	userDir := filepath.Join(u.uploadDir, userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, err
	}
	storagePath := filepath.Join(userDir, storageName)

	// 4. Save to disk
	dst, err := os.Create(storagePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(storagePath)
		return nil, err
	}

	// 5. Save to DB
	node := &file_entity.Node{
		Name:         fileHeader.Filename,
		IsDir:        false,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		Size:         fileHeader.Size,
		Extension:    ext,
		StorageName:  storageName,
		StoragePath:  storagePath,
		UserID:       userID,
		ParentID:     parentID,
	}

	if err := u.repo.Create(ctx, node); err != nil {
		os.Remove(storagePath)
		return nil, err
	}

	// 6. Log activity
	go func() {
		_ = u.activityRepo.Create(context.Background(), &file_entity.ActivityLog{
			UserID: userID,
			Action: "upload",
			FileID: &node.ID,
		})
	}()

	return node, nil
}

func (u *NodeUseCase) ListFiles(ctx context.Context, userID string, filter node_repository.FileFilter) ([]file_entity.Node, int64, error) {
	return u.repo.FindByUser(ctx, userID, filter)
}

func (u *NodeUseCase) DeleteNode(ctx context.Context, userID, id string) error {
	node, err := u.repo.FindByID(ctx, id)
	if err != nil || node.UserID != userID {
		return errors.New("file not found")
	}

	if err := u.repo.Delete(ctx, id, userID); err != nil {
		return err
	}

	if !node.IsDir {
		os.Remove(node.StoragePath)
	}

	go func() {
		_ = u.activityRepo.Create(context.Background(), &file_entity.ActivityLog{
			UserID: userID,
			Action: "delete",
			FileID: &id,
		})
	}()

	return nil
}

func (u *NodeUseCase) GetQuota(ctx context.Context, userID string) (map[string]interface{}, error) {
	used, err := u.repo.GetUsedSpace(ctx, userID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": u.maxQuota,
		"used":  used,
		"free":  u.maxQuota - used,
		"pct":   float64(used) / float64(u.maxQuota) * 100,
	}, nil
}

func (u *NodeUseCase) ToggleFavorite(ctx context.Context, userID, id string, isFavorite bool) error {
	node, err := u.repo.FindByID(ctx, id)
	if err != nil || node.UserID != userID {
		return errors.New("file not found")
	}
	return u.repo.UpdateFavorite(ctx, id, userID, isFavorite)
}

func (u *NodeUseCase) GetNode(ctx context.Context, userID, id string) (*file_entity.Node, error) {
	node, err := u.repo.FindByID(ctx, id)
	if err != nil || node.UserID != userID {
		return nil, errors.New("file not found")
	}
	return node, nil
}
