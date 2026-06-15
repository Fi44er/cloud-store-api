package file_module

import (
	file_http "github.com/Fi44er/cloud-store-api/internal/modules/files/delivery/http"
	activity_repository "github.com/Fi44er/cloud-store-api/internal/modules/files/infrastructure/repository/activity"
	node_repository "github.com/Fi44er/cloud-store-api/internal/modules/files/infrastructure/repository/node"
	file_usecase "github.com/Fi44er/cloud-store-api/internal/modules/files/usecase"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FileModule struct {
	logger *logger.Logger
	db     *gorm.DB

	nodeUseCase        *file_usecase.NodeUseCase
	activityUseCase    *file_usecase.ActivityUseCase
	nodeRepository     *node_repository.NodeRepository
	activityRepository *activity_repository.ActivityRepository
	fileHandler        *file_http.NodeHandler
	activityHandler    *file_http.ActivityHandler
}

func NewFileModule(logger *logger.Logger, db *gorm.DB) *FileModule {
	return &FileModule{
		logger: logger,
		db:     db,
	}
}

func (m *FileModule) Init() {
	m.nodeRepository = node_repository.NewNodeRepository(m.logger, m.db)
	m.activityRepository = activity_repository.NewActivityRepository(m.logger, m.db)
	
	uploadDir := "./uploads"
	maxQuota := int64(1073741824) // 1GB default

	m.nodeUseCase = file_usecase.NewNodeUseCase(m.logger, m.nodeRepository, m.activityRepository, uploadDir, maxQuota)
	m.activityUseCase = file_usecase.NewActivityUseCase(m.logger, m.activityRepository)
	
	m.fileHandler = file_http.NewNodeHandler(m.logger, m.nodeUseCase)
	m.activityHandler = file_http.NewActivityHandler(m.logger, m.activityUseCase)
}

func (m *FileModule) RegisterRoutes(router fiber.Router) {
	m.fileHandler.RegisterRoutes(router)
	m.activityHandler.RegisterRoutes(router)
}
