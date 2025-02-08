package notification

import (
	"context"
	models "eventdrivensystem/internal/models/notification"

	"gorm.io/gorm"
)

func (u *NotificationDomain) getNotificationSql(ctx context.Context, p *models.GetNotificationParam) (*models.Notification, error) {

	var (
		db     *gorm.DB
		result models.Notification
	)

	db = u.db.WithContext(ctx)
	if p.ID != "" {
		db = db.Where("id = ?", p.ID)
	}

	err := db.First(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}
