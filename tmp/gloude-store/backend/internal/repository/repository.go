// Пакет repository отвечает за взаимодействие с базой данных
// Используется GORM: https://gorm.io/docs/
package repository

import (
	"context"
	"time"

	"github.com/gloude/store/internal/models"
	"gorm.io/gorm"
)

// FileRepository — интерфейс репозитория файлов
type FileRepository interface {
	Create(ctx context.Context, file *models.File) error
	FindByID(ctx context.Context, id, userID uint) (*models.File, error)
	FindByUser(ctx context.Context, userID uint, filter models.FileFilter) ([]models.File, int64, error)
	Delete(ctx context.Context, id, userID uint) error
	GetUsedSpace(ctx context.Context, userID uint) (int64, error)
	UpdateFavorite(ctx context.Context, id, userID uint, isFavorite bool) error
}

// ActivityRepository — интерфейс репозитория активности
type ActivityRepository interface {
	Create(ctx context.Context, log *models.ActivityLog) error
	GetActivityByDays(ctx context.Context, userID uint, from time.Time) ([]models.ActivityByDay, error)
}

// UserRepository — интерфейс репозитория пользователей
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
}

// --- Реализации ---

type fileRepo struct {
	db *gorm.DB
}

type activityRepo struct {
	db *gorm.DB
}

type userRepo struct {
	db *gorm.DB
}

// NewFileRepository создает новый экземпляр репозитория файлов
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepo{db: db}
}

// NewActivityRepository создает новый экземпляр репозитория активности
func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepo{db: db}
}

// NewUserRepository создает новый экземпляр репозитория пользователей
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// ===== FileRepository =====

// Create сохраняет метаданные файла в базу данных
func (r *fileRepo) Create(ctx context.Context, file *models.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

// FindByID ищет файл по ID с проверкой принадлежности пользователю
func (r *fileRepo) FindByID(ctx context.Context, id, userID uint) (*models.File, error) {
	var file models.File
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// FindByUser возвращает список файлов пользователя с фильтрацией и пагинацией
func (r *fileRepo) FindByUser(ctx context.Context, userID uint, filter models.FileFilter) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	// Строим запрос с фильтрами
	query := r.db.WithContext(ctx).Model(&models.File{}).Where("user_id = ?", userID)

	// Фильтр по расширению
	if filter.Extension != "" {
		query = query.Where("extension = ?", filter.Extension)
	}

	// Фильтр по диапазону дат
	if !filter.DateFrom.IsZero() {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if !filter.DateTo.IsZero() {
		query = query.Where("created_at <= ?", filter.DateTo)
	}

	// Фильтр по размеру
	if filter.MinSize > 0 {
		query = query.Where("size >= ?", filter.MinSize)
	}
	if filter.MaxSize > 0 {
		query = query.Where("size <= ?", filter.MaxSize)
	}

	// Поиск по имени (ILIKE — регистронезависимо в PostgreSQL)
	if filter.Search != "" {
		query = query.Where("original_name ILIKE ?", "%"+filter.Search+"%")
	}

	// Считаем общее количество для пагинации
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Пагинация
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Сортировка по дате создания (новые первыми)
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// Delete удаляет запись о файле из базы данных
func (r *fileRepo) Delete(ctx context.Context, id, userID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.File{}).Error
}

// GetUsedSpace возвращает суммарный объем файлов пользователя в байтах
func (r *fileRepo) GetUsedSpace(ctx context.Context, userID uint) (int64, error) {
	var used int64
	err := r.db.WithContext(ctx).
		Model(&models.File{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(size), 0)").
		Scan(&used).Error
	return used, err
}

// UpdateFavorite обновляет статус избранного у файла
func (r *fileRepo) UpdateFavorite(ctx context.Context, id, userID uint, isFavorite bool) error {
	return r.db.WithContext(ctx).
		Model(&models.File{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_favorite", isFavorite).Error
}

// ===== ActivityRepository =====

// Create записывает событие активности в журнал
func (r *activityRepo) Create(ctx context.Context, log *models.ActivityLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// GetActivityByDays возвращает агрегированную активность по дням за последние 365 дней
// Используется для построения GitHub-style heatmap
func (r *activityRepo) GetActivityByDays(ctx context.Context, userID uint, from time.Time) ([]models.ActivityByDay, error) {
	var result []models.ActivityByDay

	// SQL запрос: группируем по датам и считаем количество событий
	// DATE_TRUNC обрезает timestamp до дня
	err := r.db.WithContext(ctx).
		Model(&models.ActivityLog{}).
		Select("TO_CHAR(DATE_TRUNC('day', created_at), 'YYYY-MM-DD') as date, COUNT(*) as count").
		Where("user_id = ? AND action = 'upload' AND created_at >= ?", userID, from).
		Group("DATE_TRUNC('day', created_at)").
		Order("date ASC").
		Scan(&result).Error

	return result, err
}

// ===== UserRepository =====

// Create создает нового пользователя в базе данных
func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail ищет пользователя по email
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID ищет пользователя по ID
func (r *userRepo) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
