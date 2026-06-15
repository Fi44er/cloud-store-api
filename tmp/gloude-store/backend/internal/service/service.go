// Пакет service содержит бизнес-логику приложения
// Документация os: https://pkg.go.dev/os
package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/gloude/store/internal/models"
	"github.com/gloude/store/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Кастомные ошибки сервиса
var (
	ErrQuotaExceeded   = errors.New("quota exceeded")
	ErrFileNotFound    = errors.New("file not found")
	ErrInvalidFileName = errors.New("invalid file name")
	ErrUnauthorized    = errors.New("unauthorized")
)

// StorageService — интерфейс сервиса управления файлами
type StorageService interface {
	Upload(ctx context.Context, userID uint, fileHeader *multipart.FileHeader) (*models.File, error)
	ListFiles(ctx context.Context, userID uint, filter models.FileFilter) ([]models.File, int64, error)
	GetFile(ctx context.Context, userID, fileID uint) (*models.File, error)
	DeleteFile(ctx context.Context, userID, fileID uint) error
	GetQuota(ctx context.Context, userID uint) (*models.QuotaInfo, error)
	GetActivity(ctx context.Context, userID uint) ([]models.ActivityByDay, error)
	ToggleFavorite(ctx context.Context, userID, fileID uint, isFavorite bool) error
}

// AuthService — интерфейс сервиса аутентификации
type AuthService interface {
	Register(ctx context.Context, email, username, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (*models.User, string, error)
	GetUser(ctx context.Context, userID uint) (*models.User, error)
}

// --- Реализации ---

type storageService struct {
	fileRepo     repository.FileRepository
	activityRepo repository.ActivityRepository
	userRepo     repository.UserRepository
	uploadDir    string
	maxQuota     int64
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

// NewStorageService создает новый экземпляр сервиса хранилища
func NewStorageService(
	fileRepo repository.FileRepository,
	activityRepo repository.ActivityRepository,
	userRepo repository.UserRepository,
	uploadDir string,
	maxQuota int64,
) StorageService {
	return &storageService{
		fileRepo:     fileRepo,
		activityRepo: activityRepo,
		userRepo:     userRepo,
		uploadDir:    uploadDir,
		maxQuota:     maxQuota,
	}
}

// NewAuthService создает новый экземпляр сервиса аутентификации
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

// ===== StorageService =====

// sanitizeFileName очищает имя файла от опасных символов (защита от Directory Traversal)
// Документация: https://owasp.org/www-community/attacks/Path_Traversal
func sanitizeFileName(name string) (string, error) {
	// Убираем все символы кроме букв, цифр, точек, дефисов и подчеркиваний
	re := regexp.MustCompile(`[^\w\.\-]`)
	safe := re.ReplaceAllString(filepath.Base(name), "_")

	// Проверяем что имя не начинается с точки (скрытые файлы)
	if strings.HasPrefix(safe, ".") {
		safe = "_" + safe
	}

	// Проверяем на пустое имя или только точки
	if safe == "" || safe == "." || safe == ".." {
		return "", ErrInvalidFileName
	}

	// Проверяем что нет path separators
	if strings.Contains(name, "..") || strings.ContainsAny(name, "/\\") {
		return "", ErrInvalidFileName
	}

	// Ограничиваем длину имени
	if len(safe) > 255 {
		safe = safe[:255]
	}

	return safe, nil
}

// isASCII проверяет что строка содержит только ASCII символы
func isASCII(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// Upload загружает файл и сохраняет метаданные в БД
// Документация Fiber uploads: https://docs.gofiber.io/api/ctx#formfile
func (s *storageService) Upload(ctx context.Context, userID uint, fileHeader *multipart.FileHeader) (*models.File, error) {
	// Проверяем квоту перед загрузкой
	used, err := s.fileRepo.GetUsedSpace(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get used space: %w", err)
	}

	// Получаем максимальную квоту пользователя
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	quota := user.QuotaMax
	if quota == 0 {
		quota = s.maxQuota
	}

	// Проверяем: не превысит ли загрузка доступную квоту
	if used+fileHeader.Size > quota {
		return nil, ErrQuotaExceeded
	}

	// Санитизируем оригинальное имя файла
	originalName, err := sanitizeFileName(fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	// Получаем расширение и MIME тип
	ext := strings.ToLower(filepath.Ext(originalName))
	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// Генерируем уникальное UUID имя для хранения на диске
	// Это предотвращает коллизии имен и скрывает оригинальные имена
	storageName := uuid.New().String() + ext

	// Создаем директорию пользователя: ./uploads/{user_id}/
	// os.MkdirAll: https://pkg.go.dev/os#MkdirAll
	userDir := filepath.Join(s.uploadDir, fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create user directory: %w", err)
	}

	// Полный путь к файлу на диске
	filePath := filepath.Join(userDir, storageName)

	// Открываем входящий файл
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Создаем файл на диске
	// os.Create: https://pkg.go.dev/os#Create
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Копируем содержимое файла
	if _, err := io.Copy(dst, src); err != nil {
		// Удаляем созданный файл в случае ошибки
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Сохраняем метаданные в базу данных
	file := &models.File{
		UserID:       userID,
		OriginalName: originalName,
		StorageName:  storageName,
		Extension:    ext,
		MimeType:     mimeType,
		Size:         fileHeader.Size,
		Path:         filePath,
	}

	if err := s.fileRepo.Create(ctx, file); err != nil {
		// Если не удалось сохранить в БД — удаляем файл с диска
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	// Записываем событие в журнал активности
	s.logActivity(ctx, userID, "upload", &file.ID)

	return file, nil
}

// logActivity записывает событие в журнал активности (не блокирующий)
func (s *storageService) logActivity(ctx context.Context, userID uint, action string, fileID *uint) {
	log := &models.ActivityLog{
		UserID: userID,
		Action: action,
		FileID: fileID,
	}
	// Используем горутину чтобы не блокировать основной запрос
	go func() {
		_ = s.activityRepo.Create(context.Background(), log)
	}()
}

// ListFiles возвращает список файлов пользователя с фильтрацией
func (s *storageService) ListFiles(ctx context.Context, userID uint, filter models.FileFilter) ([]models.File, int64, error) {
	return s.fileRepo.FindByUser(ctx, userID, filter)
}

// GetFile возвращает файл по ID (с проверкой принадлежности)
func (s *storageService) GetFile(ctx context.Context, userID, fileID uint) (*models.File, error) {
	file, err := s.fileRepo.FindByID(ctx, fileID, userID)
	if err != nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

// DeleteFile удаляет файл с диска и из базы данных
func (s *storageService) DeleteFile(ctx context.Context, userID, fileID uint) error {
	// Получаем метаданные файла
	file, err := s.fileRepo.FindByID(ctx, fileID, userID)
	if err != nil {
		return ErrFileNotFound
	}

	// Удаляем запись из базы данных
	if err := s.fileRepo.Delete(ctx, fileID, userID); err != nil {
		return fmt.Errorf("failed to delete file record: %w", err)
	}

	// Удаляем файл с диска
	// os.Remove: https://pkg.go.dev/os#Remove
	if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
		// Логируем ошибку но не возвращаем — запись уже удалена из БД
		fmt.Printf("warning: failed to delete file from disk: %v\n", err)
	}

	// Записываем событие в журнал активности
	s.logActivity(ctx, userID, "delete", &fileID)

	return nil
}

// GetQuota возвращает информацию о квоте пользователя
func (s *storageService) GetQuota(ctx context.Context, userID uint) (*models.QuotaInfo, error) {
	used, err := s.fileRepo.GetUsedSpace(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get used space: %w", err)
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	total := user.QuotaMax
	if total == 0 {
		total = s.maxQuota
	}

	free := total - used
	if free < 0 {
		free = 0
	}

	pct := 0.0
	if total > 0 {
		pct = float64(used) / float64(total) * 100
	}

	return &models.QuotaInfo{
		Total: total,
		Used:  used,
		Free:  free,
		Pct:   pct,
	}, nil
}

// GetActivity возвращает данные активности за последние 365 дней
func (s *storageService) GetActivity(ctx context.Context, userID uint) ([]models.ActivityByDay, error) {
	from := time.Now().AddDate(-1, 0, 0) // 365 дней назад
	return s.activityRepo.GetActivityByDays(ctx, userID, from)
}

// ToggleFavorite обновляет статус избранного файла
func (s *storageService) ToggleFavorite(ctx context.Context, userID, fileID uint, isFavorite bool) error {
	return s.fileRepo.UpdateFavorite(ctx, fileID, userID, isFavorite)
}

// ===== AuthService =====

// Register регистрирует нового пользователя
func (s *authService) Register(ctx context.Context, email, username, password string) (*models.User, error) {
	// Хешируем пароль через bcrypt
	// bcrypt: https://pkg.go.dev/golang.org/x/crypto/bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:    strings.ToLower(strings.TrimSpace(email)),
		Username: strings.TrimSpace(username),
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login проверяет учетные данные и возвращает токен JWT
func (s *authService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return nil, "", ErrUnauthorized
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", ErrUnauthorized
	}

	// Генерируем JWT токен
	token, err := generateJWT(user, s.jwtSecret)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// GetUser возвращает пользователя по ID
func (s *authService) GetUser(ctx context.Context, userID uint) (*models.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
