// Точка входа сервера Gloude Store API
// Документация Go Fiber: https://docs.gofiber.io/
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gloude/store/internal/config"
	"github.com/gloude/store/internal/handler"
	"github.com/gloude/store/internal/middleware"
	"github.com/gloude/store/internal/models"
	"github.com/gloude/store/internal/repository"
	"github.com/gloude/store/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	// Загружаем конфигурацию из переменных окружения
	cfg := config.Load()

	// Инициализируем подключение к базе данных
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Создаем директорию для загрузок если не существует
	// os.MkdirAll: https://pkg.go.dev/os#MkdirAll
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Инициализируем слои приложения
	// Репозитории (работа с БД)
	fileRepo := repository.NewFileRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Сервисы (бизнес-логика)
	storageSvc := service.NewStorageService(fileRepo, activityRepo, userRepo, cfg.UploadDir, cfg.MaxQuotaBytes)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)

	// Хендлеры (HTTP контроллеры)
	storageHandler := handler.NewStorageHandler(storageSvc)
	authHandler := handler.NewAuthHandler(authSvc)

	// Создаем Fiber приложение
	// Документация: https://docs.gofiber.io/api/fiber
	app := fiber.New(fiber.Config{
		// Максимальный размер тела запроса (1.1 ГБ — чуть больше квоты)
		BodyLimit: 1100 * 1024 * 1024,
		// Заголовок для реального IP за прокси
		ProxyHeader: fiber.HeaderXForwardedFor,
		// Кастомный обработчик ошибок
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Подключаем глобальные middleware
	// Recover: https://docs.gofiber.io/api/middleware/recover
	app.Use(recover.New())

	// Logger: https://docs.gofiber.io/api/middleware/logger
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// CORS: https://docs.gofiber.io/api/middleware/cors
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:80,http://localhost",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	// Healthcheck эндпоинт (без авторизации)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
		})
	})

	// API роуты v1
	api := app.Group("/api/v1")

	// Публичные роуты аутентификации
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	// Защищенные роуты (требуют JWT)
	authMw := middleware.AuthRequired(cfg.JWTSecret)

	// Аккаунт
	account := api.Group("/account", authMw)
	account.Get("/me", authHandler.Me)

	// Хранилище
	storage := api.Group("/storage", authMw)
	storage.Post("/upload", storageHandler.Upload)
	storage.Get("/files", storageHandler.ListFiles)
	storage.Get("/download/:file_id", storageHandler.Download)
	storage.Delete("/:file_id", storageHandler.Delete)
	storage.Get("/quota", storageHandler.GetQuota)
	storage.Get("/activity", storageHandler.GetActivity)
	storage.Patch("/:file_id/favorite", storageHandler.ToggleFavorite)

	// Запускаем сервер
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Gloude Store API server started on %s", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// initDB инициализирует подключение к PostgreSQL через GORM
// Документация GORM: https://gorm.io/docs/connecting_to_the_database.html
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// Настройка логирования GORM
	logLevel := gormLogger.Silent
	if cfg.Env == "development" {
		logLevel = gormLogger.Info
	}

	gormCfg := &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	}

	// Подключаемся с повторными попытками (для Docker Compose)
	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormCfg)
		if err == nil {
			break
		}
		log.Printf("Database connection attempt %d/10 failed, retrying in 3s...", i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after 10 attempts: %w", err)
	}

	// Настраиваем пул соединений
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Автомиграция — создание/обновление таблиц
	// GORM AutoMigrate: https://gorm.io/docs/migration.html
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.File{},
		&models.ActivityLog{},
	); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Database migrations completed successfully")

	return db, nil
}
