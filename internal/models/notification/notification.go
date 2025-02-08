package notification

import (
	"time"

	"github.com/go-openapi/strfmt"
)

// Notification model
type Notification struct {
	ID        strfmt.UUID4 `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id" json:"id"`
	UserID    strfmt.UUID4 `gorm:"type:uuid;not null;column:user_id" json:"user_id"`
	Type      string       `gorm:"not null;column:type" json:"type"`
	Message   string       `gorm:"not null;column:message" json:"message"`
	Status    string       `gorm:"not null;column:status" json:"status"`
	CreatedAt time.Time    `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	DeletedAt *time.Time   `gorm:"column:deleted_at" json:"deleted_at"`
}

func (n *Notification) TableName() string {
	return "notifications"
}

type GetNotificationParam struct {
	ID strfmt.UUID4 `json:"id"`
}
