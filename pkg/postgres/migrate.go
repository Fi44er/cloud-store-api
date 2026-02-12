package postgres

import (
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB, trigger bool, log *logger.Logger) error {

	if trigger {
		log.Info("📦 Migrating database...")
		models := []any{}

		log.Info("📦 Creating types...")

		db.Exec("CREATE TYPE file_status AS ENUM ('temporary', 'permanent')")
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

		if err := db.AutoMigrate(models...); err != nil {
			log.Errorf("✖ Failed to migrate database: %v", err)
			return err
		}
	}

	log.Info("✅ Database connection successfully")
	return nil
}
