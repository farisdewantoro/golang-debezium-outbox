// kafka/kafka.go
package kafka

import (
	"context"
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/logger"
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaClient manages the Kafka producer and consumers.
type KafkaClient struct {
	cfg       *configs.AppConfig
	log       logger.Logger
	Producer  *kafka.Producer
	Consumers *KafkaConsumerManager
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

// KafkaConsumerManager handles multiple consumers.
type KafkaConsumerManager struct {
	mu        sync.Mutex
	consumers map[string]*KafkaWorkerPool
}

type KafkaWorkerPool struct {
	Consumer *kafka.Consumer
	Option   configs.KafkaTopicOption
}

// NewKafkaClient initializes a Kafka client with a producer and consumer manager.
func NewKafkaClient(cfg *configs.AppConfig, log logger.Logger) (*KafkaClient, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create Kafka producer
	producerConfig := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Address,
	}

	producer, err := kafka.NewProducer(producerConfig)
	if err != nil {
		cancel()
		return nil, err
	}

	client := &KafkaClient{
		Producer: producer,
		Consumers: &KafkaConsumerManager{
			consumers: make(map[string]*KafkaWorkerPool),
		},
		cfg:    cfg,
		log:    log,
		ctx:    ctx,
		cancel: cancel,
		wg:     sync.WaitGroup{},
	}

	return client, nil
}

// PublishMessage asynchronously sends a message to Kafka.
func (k *KafkaClient) PublishMessage(topic string, message []byte) error {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
	}

	deliveryChan := make(chan kafka.Event)
	err := k.Producer.Produce(msg, deliveryChan)
	if err != nil {
		k.log.Errorf("Failed to produce message: %v", err)
		return err
	}

	// Handle delivery response asynchronously
	go func() {
		defer close(deliveryChan)
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			k.log.Errorf("Failed to produce message: %v", m.TopicPartition.Error)
		} else {
			k.log.Infof("Message sent to topic %s [partition %d] at offset %d",
				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
		}
	}()

	return nil
}

// RegisterConsumer creates a new Kafka consumer and subscribes to a topic.
func (k *KafkaClient) RegisterConsumer(opt configs.KafkaTopicOption) (*KafkaWorkerPool, error) {
	k.Consumers.mu.Lock()
	defer k.Consumers.mu.Unlock()

	// Check if consumer already exists
	if consumer, exists := k.Consumers.consumers[opt.Topic]; exists {
		return consumer, nil
	}

	groupID := &k.cfg.Kafka.GroupID

	if opt.GroupID != nil {
		groupID = opt.GroupID
	}

	consumerConfig := &kafka.ConfigMap{
		"bootstrap.servers":  k.cfg.Kafka.Address,
		"group.id":           *groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	}

	consumer, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		return nil, err
	}

	if err := consumer.Subscribe(opt.Topic, nil); err != nil {
		return nil, err
	}

	kafkaWorkerPool := &KafkaWorkerPool{
		Consumer: consumer,
		Option:   opt,
	}
	k.Consumers.consumers[opt.Topic] = kafkaWorkerPool
	k.log.Infof("Registered consumer for topic: %s group:%s", opt.Topic, *groupID)
	return kafkaWorkerPool, nil
}

func (k *KafkaClient) RegisterConsumerRegex(opt configs.KafkaTopicOption) (*KafkaWorkerPool, error) {
	k.Consumers.mu.Lock()
	defer k.Consumers.mu.Unlock()

	// Check if consumer already exists
	if consumer, exists := k.Consumers.consumers[opt.Topic]; exists {
		return consumer, nil
	}

	groupID := &k.cfg.Kafka.GroupID

	if opt.GroupID != nil {
		groupID = opt.GroupID
	}

	consumerConfig := &kafka.ConfigMap{
		"bootstrap.servers":  k.cfg.Kafka.Address,
		"group.id":           *groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	}

	consumer, err := kafka.NewConsumer(consumerConfig)
	if err != nil {
		return nil, err
	}

	if err := consumer.SubscribeTopics([]string{opt.Topic}, nil); err != nil {
		return nil, err
	}

	kafkaWorkerPool := &KafkaWorkerPool{
		Consumer: consumer,
		Option:   opt,
	}
	k.Consumers.consumers[opt.Topic] = kafkaWorkerPool
	k.log.Infof("Registered consumer regex for topic: %s group:%s", opt.Topic, *groupID)
	return kafkaWorkerPool, nil
}

// ConsumeMessages starts consuming messages from a topic.
func (k *KafkaClient) ConsumeMessages(topic string, handler func(context.Context, *kafka.Message) error) error {
	worker, exists := k.Consumers.consumers[topic]
	if !exists {
		return fmt.Errorf("consumer for topic %s not registered", topic)
	}

	k.log.Infof("Consuming messages from topic: %s", topic)

	// Worker pool
	maxConcurrency := make(chan struct{}, worker.Option.MaxConcurrency)

	k.wg.Add(1)
	go func() {
		defer k.wg.Done()
		for {
			select {
			case <-k.ctx.Done():
				k.log.Info("Stopping Kafka consumer:", topic)
				return
			default:
				msg, err := worker.Consumer.ReadMessage(worker.Option.PoolDuration)

				if err != nil {
					kafkaErr, ok := err.(kafka.Error)

					if ok && kafkaErr.IsRetriable() {
						k.log.Warn("Temporary error reading message Retrying...", err)
						time.Sleep(2 * time.Second) // Simple backoff
						continue
					} else if ok && kafkaErr.IsTimeout() {
						continue
					} else {
						k.log.Errorf("Fatal error reading message: %v. Exiting consumer loop.", err)
						return
					}

				}

				// Limit concurrency
				maxConcurrency <- struct{}{}
				go func() {
					defer func() {
						<-maxConcurrency
					}()

					k.processMessageWithRetry(msg, worker.Option, handler)
				}()

			}
		}
	}()

	return nil
}

// processMessageWithRetry retries message processing with exponential backoff.
func (k *KafkaClient) processMessageWithRetry(msg *kafka.Message, opt configs.KafkaTopicOption, handler func(context.Context, *kafka.Message) error) {
	retryPolicy := backoff.NewExponentialBackOff()
	retryPolicy.InitialInterval = opt.RetryInitIntervalDuration
	retryPolicy.MaxInterval = opt.RetryMaxIntervalDuration
	retryPolicy.MaxElapsedTime = opt.RetryMaxElapsedTimeDuration
	retryPolicy.RandomizationFactor = 0.5 // Jitter

	operation := func() error {
		select {
		case <-k.ctx.Done(): // stop processing if server is shutting down
			return context.Canceled
		default:
			ctx, cancel := context.WithTimeout(context.Background(), opt.MaxTimeoutDuration)
			defer cancel()

			err := handler(ctx, msg)
			if err != nil {
				k.log.Errorf("Retrying due to error: %v", err)
			}
			return err
		}
	}

	err := backoff.Retry(operation, retryPolicy)
	if err != nil {
		k.log.Errorf("Max retries reached, sending to DLQ: %s", string(msg.Value))
		_ = k.PublishMessage(k.cfg.Kafka.DLQTopic, msg.Value)
	}
}

// Close gracefully shuts down the Kafka client.
func (k *KafkaClient) Close() {
	k.log.Info("Shutting down Kafka client...")

	// Signal cancellation
	k.cancel()

	// Wait for ongoing consumers to finish
	k.wg.Wait()

	// Close all consumers
	k.Consumers.mu.Lock()
	for topic, worker := range k.Consumers.consumers {
		k.log.Infof("Closing consumer for topic: %s", topic)
		worker.Consumer.Close()
	}
	k.Consumers.mu.Unlock()

	// Close Kafka producer
	k.Producer.Close()

	k.log.Info("Kafka client shut down successfully.")
}
