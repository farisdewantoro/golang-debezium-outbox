package kafka

import (
	"context"
	"eventdrivensystem/configs"
	"eventdrivensystem/internal/usecase"
	kafkaPkg "eventdrivensystem/pkg/kafka"
	"eventdrivensystem/pkg/logger"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumerHandler struct {
	cfg         *configs.AppConfig
	log         logger.Logger
	kafkaClient *kafkaPkg.KafkaClient
	uc          *usecase.Usecase
}

func NewKafkaConsumerHandler(cfg *configs.AppConfig, log logger.Logger, kafkaClient *kafkaPkg.KafkaClient, uc *usecase.Usecase) *KafkaConsumerHandler {
	return &KafkaConsumerHandler{
		cfg:         cfg,
		log:         log,
		kafkaClient: kafkaClient,
		uc:          uc,
	}
}

type HandlerFunc func(ctx context.Context, msg *kafka.Message) error

func (k *KafkaConsumerHandler) Handle(opt configs.KafkaTopicOption, handler HandlerFunc) {
	if opt.UseRegex {
		_, err := k.kafkaClient.RegisterConsumerRegex(opt)
		if err != nil {
			log.Fatalf("failed to register consumer: %v", err)
		}
	} else {
		_, err := k.kafkaClient.RegisterConsumer(opt)
		if err != nil {
			log.Fatalf("failed to register consumer: %v", err)
		}
	}

	k.kafkaClient.ConsumeMessages(opt.Topic, handler)
}

func (k *KafkaConsumerHandler) StartConsumers() {
	if k.cfg.Kafka.Options.OutboxEvent.Enable {
		k.Handle(k.cfg.Kafka.Options.OutboxEvent, k.ProcessOutboxEvent)
	}

	if k.cfg.Kafka.Options.NotifUserRegistrationEvent.Enable {
		k.Handle(k.cfg.Kafka.Options.NotifUserRegistrationEvent, k.ProcessNotifUserRegistrationEvent)
	}
}
