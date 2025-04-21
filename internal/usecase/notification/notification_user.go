package notification

import (
	"context"
	models "eventdrivensystem/internal/models/notification"
	"eventdrivensystem/pkg/errors"

	"github.com/go-openapi/strfmt"
)

func (u *NotificationUsecase) SendNotificationUserRegistration(ctx context.Context, id strfmt.UUID4) error {
	u.log.InfoWithContext(ctx, "Sending notification to user %s", id)
	existingNotif, err := u.notificationDomain.GetNotification(ctx, &models.GetNotificationParam{ID: id})

	if err != nil {
		u.log.ErrorWithContext(ctx, err)
		return errors.ErrSQLGet
	}

	if existingNotif.Status == models.NotificationStatusSent {
		u.log.WarnWithContext(ctx, "Notification already sent to user: ", existingNotif.UserID)
		return nil
	}

	// Create new notification object with exact constant values
	notif := &models.Notification{
		ID:      existingNotif.ID,
		UserID:  existingNotif.UserID,
		Type:    models.NotificationTypeUserRegistration,
		Message: models.NotificationMessageUserRegistration,
		Status:  models.NotificationStatusSent,
	}

	updatedNotif, err := u.notificationDomain.UpdateNotification(ctx, notif)
	if err != nil {
		u.log.ErrorWithContext(ctx, err)
		return errors.ErrSQLCreate
	}

	if updatedNotif == nil {
		return errors.ErrSQLCreate
	}

	u.log.InfofWithContext(ctx, "Notification sent to user: %v", notif.UserID)

	return nil
}
