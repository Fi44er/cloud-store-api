package node_repository

import (
	file_entity "github.com/Fi44er/cloud-store-api/internal/modules/files/entity"
	file_model "github.com/Fi44er/cloud-store-api/internal/modules/files/infrastructure/repository/model"
)

type converter struct {
}

func newConverter() *converter {
	return &converter{}
}

func (c *converter) toModel(node *file_entity.Node) *file_model.Node {
	return &file_model.Node{
		ID:       node.ID,
		Name:     node.Name,
		IsDir:    node.IsDir,
		MimeType: node.MimeType,
		Size:     node.Size,
		UserID:   node.UserID,
		ParentID: node.ParentID,
	}
}

func (c *converter) toEntity(node *file_model.Node) *file_entity.Node {
	return &file_entity.Node{
		ID:       node.ID,
		Name:     node.Name,
		IsDir:    node.IsDir,
		MimeType: node.MimeType,
		Size:     node.Size,
		UserID:   node.UserID,
		ParentID: node.ParentID,
	}
}
