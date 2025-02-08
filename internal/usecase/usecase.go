package usecase

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/internal/domain"
	"eventdrivensystem/internal/usecase/notification"
	"eventdrivensystem/internal/usecase/user"
	"eventdrivensystem/pkg/logger"
)

type Usecase struct {
	User         user.UserUsecaseHandler
	Notification notification.NotificationUsecaseHandler
}

func NewUsecase(
	cfg *configs.AppConfig,
	log logger.Logger,
	dom *domain.Domain,
) *Usecase {
	return &Usecase{
		User:         user.NewUserUsecase(cfg, log, dom),
		Notification: notification.NewNotificationUsecase(cfg, log, dom),
	}
}
