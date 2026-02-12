package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fi44er/cloud-store-api/internal/config"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/Fi44er/cloud-store-api/pkg/postgres"
	"github.com/Fi44er/cloud-store-api/pkg/postgres/uow"
	"github.com/Fi44er/cloud-store-api/pkg/process_manager"
	redisConnect "github.com/Fi44er/cloud-store-api/pkg/redis"
	"github.com/Fi44er/cloud-store-api/pkg/session"
	sessionstore "github.com/Fi44er/cloud-store-api/pkg/session/store"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/swag/example/basic/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	app *fiber.App

	config     *config.Config
	logger     *logger.Logger
	validator  *validator.Validate
	httpConfig config.HTTPConfig

	db          *gorm.DB
	redisClient *redis.Client

	redisManager   redisConnect.IRedisManager
	sessionManager *session.SessionManager
	processManager process_manager.IProcessManager
	uow            uow.Uow

	migrate   bool
	redisMode int
}

func (app *App) initConfig() error {
	if app.config == nil {
		config, err := config.LoadConfig(".")
		if err != nil {
			return fmt.Errorf("✖ Failed to load config: %s", err.Error())
		}
		app.config = config
	}

	err := config.Load(".env")
	if err != nil {
		return fmt.Errorf("✖ Failed to load config: %s", err.Error())
	}

	return nil
}

func (app *App) initDb() error {
	if app.db == nil {
		db, err := postgres.ConnectDb(app.config.DatabaseURL, app.logger)
		if err != nil {
			return err
		}
		app.db = db
		app.uow = uow.New(app.db)

		if err := postgres.Migrate(db, app.migrate, app.logger); err != nil {
			return fmt.Errorf("✖ Failed to migrate database: %s", err.Error())
		}
	}

	return nil
}

func (app *App) initRedis() error {
	if app.redisManager == nil {
		client, err := redisConnect.Connect(app.config.RedisURL, app.logger)
		if err != nil {
			app.logger.Errorf("Failed to connect to Redis: %v", err)
			return nil
		}

		app.redisManager = redisConnect.NewRedisManger(client)
		app.redisClient = client

		if err := redisConnect.FlushRedisCache(client, app.redisMode, app.logger); err != nil {
			err = fmt.Errorf("✖ Failed to flush redis cache: %v", err)
			app.logger.Errorf("%s", err.Error())
			return err
		}
	}
	return nil
}

func (app *App) initLogger() error {
	if app.logger == nil {
		app.logger = logger.NewLogger()
	}
	return nil
}

func (app *App) initValidator() error {
	if app.validator == nil {
		app.validator = validator.New()
	}
	return nil
}

func (app *App) initSessionManager() error {
	if app.sessionManager == nil {
		app.sessionManager = session.NewSessionManager(
			sessionstore.NewRedisSessionStore(app.redisClient),
			30*time.Minute,
			2*time.Hour,
			12*time.Hour,
			"session",
		)
	}

	return nil
}

func (app *App) runHttpServerWithShutdown() error {
	if app.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			app.logger.Errorf("✖ Failed to load config: %s", err.Error())
			return fmt.Errorf("✖ Failed to load config: %v", err)
		}
		app.httpConfig = cfg
	}

	serverErr := make(chan error, 1)

	go func() {
		app.logger.Infof("🌐 Server is running on %s", app.httpConfig.Address())
		app.logger.Info("✅ Server started successfully")
		if err := app.app.Listen(app.httpConfig.Address()); err != nil {
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		app.logger.Errorf("✖ Server error: %s", err.Error())
		app.stopBackgroundProcesses()
		return err
	case sig := <-quit:
		app.logger.Infof("Received signal: %v. Shutting down...", sig)
		app.stopBackgroundProcesses()
		app.logger.Info("✅ Application stopped gracefully")
		return nil
	}
}

func (app *App) stopBackgroundProcesses() {
	if app.processManager != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := app.processManager.StopAll(ctx); err != nil {
			app.logger.Errorf("Error stopping background processes: %v", err)
		}
	}
}

func (app *App) initProcessManager() error {
	if app.processManager == nil {
		app.processManager = process_manager.NewProcessManager(app.logger)
	}
	return nil
}

func (app *App) registerBackgroundProcesses() error {
	if app.processManager == nil {
		return fmt.Errorf("process manager is not initialized")
	}

	return nil
}

func (app *App) initRouter() error {
	docs.SwaggerInfo.Host = app.config.ExternalHost
	app.app.Get("/swagger/*", swagger.HandlerDefault)

	// api := app.app.Group("/api")

	return nil
}
