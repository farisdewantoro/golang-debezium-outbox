package notification

import (
	"context"
	models "eventdrivensystem/internal/models/notification"
)

func (u *NotificationDomain) GetNotification(ctx context.Context, p *models.GetNotificationParam) (*models.Notification, error) {
	return u.getNotificationSql(ctx, p)
}
