package file_model

import (
	"time"
)

type ActivityLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Action    string    `gorm:"not null;size:50" json:"action"` // "upload" or "delete"
	FileID    *string   `gorm:"type:uuid" json:"file_id,omitempty"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
