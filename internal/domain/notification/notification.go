package notification

import (
	"context"
	"eventdrivensystem/configs"
	models "eventdrivensystem/internal/models/notification"
	"eventdrivensystem/pkg/logger"
	"eventdrivensystem/pkg/util"

	"gorm.io/gorm"
)

type NotificationDomain struct {
	cfg *configs.AppConfig
	db  *gorm.DB
	log logger.Logger
}

type NotificationDomainHandler interface {
	BeginTx(ctx context.Context) *gorm.DB
	CreateNotification(ctx context.Context, p *models.Notification, opts ...util.DbOptions) (*models.Notification, error)
	UpdateNotification(ctx context.Context, p *models.Notification, opts ...util.DbOptions) (*models.Notification, error)
	GetNotification(ctx context.Context, p *models.GetNotificationParam) (*models.Notification, error)
}

func NewNotificationDomain(cfg *configs.AppConfig, log logger.Logger, db *gorm.DB) NotificationDomainHandler {
	return &NotificationDomain{
		cfg: cfg,
		db:  db,
		log: log,
	}
}
