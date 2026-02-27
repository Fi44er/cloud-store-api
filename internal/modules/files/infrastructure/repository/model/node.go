package file_model

import (
	"time"

	"gorm.io/gorm"
)

type Node struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Мягкое удаление (корзина)

	Name     string `gorm:"not null" json:"name"`   // Имя файла: "отпуск.jpg"
	IsDir    bool   `gorm:"not null" json:"is_dir"` // true = папка, false = файл
	MimeType string `json:"mime_type"`              // image/jpeg, application/pdf
	Size     int64  `json:"size"`                   // Размер в байтах

	// Физический путь на диске (только для файлов).
	// Например: "uploads/2023/05/a1-b2-c3..."
	StoragePath string `json:"-"`

	// Владелец
	UserID string `gorm:"type:uuid;index" json:"user_id"`
	// User   User      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`

	// Иерархия (ссылка на саму себя)
	// Если ParentID == nil, значит файл/папка лежит в корне
	ParentID *string `gorm:"type:uuid;index" json:"parent_id"`
	Children []Node  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
