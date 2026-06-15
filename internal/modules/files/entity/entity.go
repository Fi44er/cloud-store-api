package file_entity

import (
	"time"
)

type Node struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	Name         string
	IsDir        bool
	MimeType     string
	Size         int64
	Extension    string
	StorageName  string
	StoragePath  string
	IsFavorite   bool

	UserID string

	ParentID *string
	Children []Node
}
