package notification

import (
	"context"
	models "eventdrivensystem/internal/models/notification"
	"eventdrivensystem/pkg/errors"

	"github.com/go-openapi/strfmt"
)

func (u *NotificationUsecase) SendNotificationUserRegistration(ctx context.Context, id strfmt.UUID4) error {
	u.log.InfoWithContext(ctx, "Sending notification to user %s", id)
	notif, err := u.notificationDomain.GetNotification(ctx, &models.GetNotificationParam{ID: id})

	if err != nil {
		u.log.ErrorWithContext(ctx, err)
		return errors.ErrSQLGet
	}

	if notif.Status == models.NotificationStatusSent {
		u.log.WarnWithContext(ctx, "Notification already sent to user: ", notif.UserID)
		return nil
	}

	notif.Status = models.NotificationStatusSent

	_, err = u.notificationDomain.UpdateNotification(ctx, notif)

	if err != nil {
		u.log.ErrorWithContext(ctx, err)
		return errors.ErrSQLCreate
	}

	u.log.InfofWithContext(ctx, "Notification sent to user: %v", notif.UserID)

	return nil
}
