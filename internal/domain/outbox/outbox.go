package outbox

import (
	"context"
	"eventdrivensystem/configs"
	models "eventdrivensystem/internal/models/outbox"
	"eventdrivensystem/pkg/logger"
	"eventdrivensystem/pkg/util"

	"gorm.io/gorm"
)

type OutboxDomain struct {
	cfg *configs.AppConfig
	db  *gorm.DB
	log logger.Logger
}

type OutboxDomainHandler interface {
	CreateOutbox(ctx context.Context, Outbox *models.OutboxEvent, opts ...util.DbOptions) error
}

func NewOutboxDomain(cfg *configs.AppConfig, log logger.Logger, db *gorm.DB) OutboxDomainHandler {
	return &OutboxDomain{
		cfg: cfg,
		db:  db,
		log: log,
	}
}
