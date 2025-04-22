package cmd

import (
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/cache"
	"eventdrivensystem/pkg/databases"
	"eventdrivensystem/pkg/kafka"
	"eventdrivensystem/pkg/logger"
	"log"
	"runtime"
	"sync"

	goValidator "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	appDependency         *AppDependency
	appDependencySyncOnce sync.Once
)

type AppDependency struct {
	db          *gorm.DB
	cfg         *configs.AppConfig
	log         logger.Logger
	validator   *goValidator.Validate
	kafkaClient *kafka.KafkaClient
	redisClient *cache.RedisClient
}

func GetAppDependency() *AppDependency {
	appDependencySyncOnce.Do(func() {
		appDependency = NewAppDependency()
	})
	return appDependency
}

func NewAppDependency() *AppDependency {
	cfg := configs.Get()
	db, err := databases.NewSqlDb(cfg)
	lgOptions := logger.Options{
		Output:    logger.OutputStdout,
		Formatter: logger.FormatJSON,
		Level:     logger.LevelInfo,
		DefaultFields: map[string]string{
			"app.name":    cfg.Meta.Name,
			"app.runtime": runtime.Version(),
		},
	}

	lg := logger.Init(lgOptions)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	kafkaClient, err := kafka.NewKafkaClient(cfg, lg)
	if err != nil {
		log.Fatalf("failed to connect to kafka: %v", err)
	}

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return &AppDependency{
		db:          db,
		cfg:         cfg,
		log:         lg,
		validator:   goValidator.New(),
		kafkaClient: kafkaClient,
		redisClient: redisClient,
	}
}
