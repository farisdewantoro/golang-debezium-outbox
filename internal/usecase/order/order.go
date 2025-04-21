package order

import (
	"context"
	"eventdrivensystem/internal/domain"
	"eventdrivensystem/internal/domain/notification"
	"eventdrivensystem/internal/domain/order"
	"eventdrivensystem/internal/domain/outbox"
	orderModels "eventdrivensystem/internal/models/order"
	"eventdrivensystem/pkg/logger"
)

type OrderUsecase struct {
	domain             *domain.Domain
	log                logger.Logger
	orderDomain        order.OrderDomainHandler
	notificationDomain notification.NotificationDomainHandler
	outboxDomain       outbox.OutboxDomainHandler
}

type OrderUsecaseHandler interface {
	CreateOrder(ctx context.Context, param *orderModels.CreateOrderParam) error
}

func NewOrderUsecase(domain *domain.Domain, log logger.Logger) OrderUsecaseHandler {
	return &OrderUsecase{
		domain:             domain,
		log:                log,
		orderDomain:        domain.Order,
		notificationDomain: domain.Notification,
		outboxDomain:       domain.Outbox,
	}
}
