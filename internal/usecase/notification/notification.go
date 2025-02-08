package notification

import (
	"context"
	"eventdrivensystem/configs"
	"eventdrivensystem/internal/domain"
	"eventdrivensystem/internal/domain/notification"
	"eventdrivensystem/internal/domain/outbox"
	"eventdrivensystem/pkg/logger"

	"github.com/go-openapi/strfmt"
)

type NotificationUsecase struct {
	cfg *configs.AppConfig
	log logger.Logger

	outboxDomain       outbox.OutboxDomainHandler
	notificationDomain notification.NotificationDomainHandler
}

type NotificationUsecaseHandler interface {
	SendNotificationUserRegistration(ctx context.Context, id strfmt.UUID4) error
}

func NewNotificationUsecase(
	cfg *configs.AppConfig,
	log logger.Logger,
	dom *domain.Domain,
) NotificationUsecaseHandler {
	return &NotificationUsecase{
		cfg:                cfg,
		log:                log,
		outboxDomain:       dom.Outbox,
		notificationDomain: dom.Notification,
	}
}
