package file_module

import (
	file_http "github.com/Fi44er/cloud-store-api/internal/modules/files/delivery/http"
	node_repository "github.com/Fi44er/cloud-store-api/internal/modules/files/infrastructure/repository/node"
	file_usecase "github.com/Fi44er/cloud-store-api/internal/modules/files/usecase"
	"github.com/Fi44er/cloud-store-api/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FileModule struct {
	logger *logger.Logger
	db     *gorm.DB

	nodeUseCase    *file_usecase.NodeUseCase
	nodeRepository *node_repository.NodeRepository
	fileHandler    *file_http.NodeHandler
}

func NewFileModule(logger *logger.Logger, db *gorm.DB) *FileModule {
	return &FileModule{
		logger: logger,
		db:     db,
	}
}

func (m *FileModule) Init() {
	m.nodeRepository = node_repository.NewNodeRepository(m.logger, m.db)
	m.nodeUseCase = file_usecase.NewNodeUseCase(m.logger, m.nodeRepository)
	m.fileHandler = file_http.NewNodeHandler(m.logger, m.nodeUseCase)
}

func (m *FileModule) RegisterRoutes(router fiber.Router) {
	m.fileHandler.RegisterRoutes(router)
}
