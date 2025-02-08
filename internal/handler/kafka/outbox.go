package kafka

import (
	"context"
	"encoding/json"
	models "eventdrivensystem/internal/models/kafka"
	notifModels "eventdrivensystem/internal/models/notification"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgtype"
)

func (h *KafkaConsumerHandler) ProcessOutboxEvent(ctx context.Context, msg *kafka.Message) error {
	h.log.InfoWithContext(ctx, "Processing outbox event")
	dbzMsg := models.DebeziumMessage{}
	if err := json.Unmarshal(msg.Value, &dbzMsg); err != nil {
		h.log.ErrorfWithContext(ctx, "Failed to unmarshal outbox event: %v", err)
		return err
	}

	outboxEvent := dbzMsg.Payload

	h.log.InfofWithContext(ctx, "Republishing outbox_event from topic %s with event type %s and id %s", msg.TopicPartition.Topic, outboxEvent.EventType, outboxEvent.AggregateID)

	var jsonbPayload pgtype.JSONB
	if err := jsonbPayload.Set(dbzMsg.Payload.Payload); err != nil {
		h.log.ErrorfWithContext(ctx, "Failed to convert payload to JSONB: %v", err)
		return err
	}

	return h.kafkaClient.PublishMessage(outboxEvent.EventType, jsonbPayload.Bytes)
}

func (h *KafkaConsumerHandler) ProcessNotifUserRegistrationEvent(ctx context.Context, msg *kafka.Message) error {

	notif := notifModels.Notification{}
	if err := json.Unmarshal(msg.Value, &notif); err != nil {
		h.log.ErrorfWithContext(ctx, "Failed to unmarshal outbox event: %v", err)
		return err
	}

	return h.uc.Notification.SendNotificationUserRegistration(ctx, notif.ID)
}
