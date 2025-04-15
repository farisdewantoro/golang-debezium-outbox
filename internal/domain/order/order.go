package order

import (
	"context"
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/logger"

	"gorm.io/gorm"
)

type OrderDomain struct {
	cfg *configs.AppConfig
	db  *gorm.DB
	log logger.Logger
}

type OrderDomainHandler interface {
	BeginTx(ctx context.Context) *gorm.DB
	OrderDomainWriter
}

func NewOrderDomain(cfg *configs.AppConfig, log logger.Logger, db *gorm.DB) OrderDomainHandler {
	return &OrderDomain{
		cfg: cfg,
		db:  db,
		log: log,
	}
}
