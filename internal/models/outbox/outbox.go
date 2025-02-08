package models

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgtype"
)

type OutboxEvent struct {
	ID            strfmt.UUID4  `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Payload       *pgtype.JSONB `json:"payload" gorm:"type:jsonb;default:'{}';not null"`
	AggregateType string        `json:"aggregate_type" gorm:"column:aggregate_type"`
	EventType     string        `json:"event_type" gorm:"column:event_type"`
	AggregateID   string        `json:"aggregate_id" gorm:"column:aggregate_id"`
	CreatedAt     time.Time     `json:"created_at" gorm:"column:created_at"`
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}
