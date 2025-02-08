package models

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgtype"
)

type OutboxEvent struct {
	ID            strfmt.UUID4  `json:"id" `
	Payload       *pgtype.JSONB `json:"payload" `
	AggregateType string        `json:"aggregate_type" `
	EventType     string        `json:"event_type" `
	AggregateID   string        `json:"aggregate_id" `
	CreatedAt     time.Time     `json:"created_at" `
}

type DebeziumMessage struct {
	Payload OutboxEvent `json:"payload"` // Extract only "payload"
}
