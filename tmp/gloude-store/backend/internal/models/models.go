// Пакет models содержит структуры данных (модели GORM)
package models

import (
	"time"
)

// User — модель пользователя
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Username  string    `gorm:"uniqueIndex;not null;size:100" json:"username"`
	Password  string    `gorm:"not null" json:"-"` // Хеш пароля, не возвращается в JSON
	QuotaMax  int64     `gorm:"default:1073741824" json:"quota_max"` // Максимальная квота в байтах (1 ГБ)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Связи
	Files         []File         `gorm:"foreignKey:UserID" json:"-"`
	ActivityLogs  []ActivityLog  `gorm:"foreignKey:UserID" json:"-"`
}

// File — модель файла, хранит метаданные о загруженном файле
// Документация GORM: https://gorm.io/docs/models.html
type File struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	OriginalName string    `gorm:"not null;size:500" json:"original_name"`   // Оригинальное имя файла от пользователя
	StorageName  string    `gorm:"not null;size:500;uniqueIndex" json:"storage_name"` // UUID имя на диске
	Extension    string    `gorm:"size:50;index" json:"extension"`            // Расширение файла (.jpg, .pdf и т.д.)
	MimeType     string    `gorm:"size:200" json:"mime_type"`                 // MIME тип (image/jpeg, application/pdf)
	Size         int64     `gorm:"not null" json:"size"`                      // Размер файла в байтах
	Path         string    `gorm:"not null;size:1000" json:"path"`            // Полный путь к файлу на диске
	IsFavorite   bool      `gorm:"default:false" json:"is_favorite"`          // Избранный файл
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Связи
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// ActivityLog — журнал действий пользователя (для GitHub-style heatmap)
type ActivityLog struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Action    string    `gorm:"not null;size:50" json:"action"` // "upload" или "delete"
	FileID    *uint     `json:"file_id,omitempty"`              // Nullable — файл мог быть удален
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// ActivityByDay — агрегированная активность за один день (для heatmap)
type ActivityByDay struct {
	Date  string `json:"date"`  // Дата в формате YYYY-MM-DD
	Count int    `json:"count"` // Количество файлов за этот день
}

// QuotaInfo — информация о квоте пользователя
type QuotaInfo struct {
	Total int64   `json:"total"`      // Всего байт (квота)
	Used  int64   `json:"used"`       // Использовано байт
	Free  int64   `json:"free"`       // Свободно байт
	Pct   float64 `json:"percentage"` // Процент использования
}

// FileFilter — фильтры для списка файлов
type FileFilter struct {
	Extension string    `query:"extension"` // Фильтр по расширению
	DateFrom  time.Time `query:"date_from"`  // Начало диапазона дат
	DateTo    time.Time `query:"date_to"`    // Конец диапазона дат
	MinSize   int64     `query:"min_size"`   // Минимальный размер в байтах
	MaxSize   int64     `query:"max_size"`   // Максимальный размер в байтах
	Search    string    `query:"search"`     // Поиск по имени
	Page      int       `query:"page"`       // Страница (для пагинации)
	Limit     int       `query:"limit"`      // Записей на страницу
}
