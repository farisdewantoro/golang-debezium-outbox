package domain

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/internal/domain/notification"
	"eventdrivensystem/internal/domain/order"
	"eventdrivensystem/internal/domain/outbox"
	"eventdrivensystem/internal/domain/user"
	"eventdrivensystem/pkg/logger"

	"gorm.io/gorm"
)

type Domain struct {
	User         user.UserDomainHandler
	Outbox       outbox.OutboxDomainHandler
	Notification notification.NotificationDomainHandler
	Order        order.OrderDomainHandler
}

func NewDomain(cfg *configs.AppConfig,
	db *gorm.DB,
	log logger.Logger) *Domain {
	return &Domain{
		User:   user.NewUserDomain(cfg, log, db),
		Outbox: outbox.NewOutboxDomain(cfg, log, db),
		Notification: notification.NewNotificationDomain(
			cfg,
			log,
			db,
		),
		Order: order.NewOrderDomain(cfg, log, db),
	}
}
