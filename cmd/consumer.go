package cmd

import (
	"eventdrivensystem/internal/domain"
	"eventdrivensystem/internal/handler/kafka"
	"eventdrivensystem/internal/usecase"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var kafkaConsumerCmd = &cobra.Command{
	Use:   "kafka-consumer",
	Short: "Runs the kafka consumer",
	Run: func(cmd *cobra.Command, args []string) {
		RunKafkaConsumer()
	},
}

func RunKafkaConsumer() {
	fmt.Println("Starting kafka consumer...")
	dp := GetAppDependency()

	dom := domain.NewDomain(dp.cfg, dp.db, dp.log)
	uc := usecase.NewUsecase(dp.cfg, dp.log, dom)

	connector := NewConnectorManager(dp.cfg, dp.log)

	if err := connector.RegisterOutboxConnector(); err != nil {
		dp.log.Fatal(err)
		return
	}

	kafkaHandler := kafka.NewKafkaConsumerHandler(dp.cfg, dp.log, dp.kafkaClient, uc)

	kafkaHandler.StartConsumers()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	dp.log.Info("shut down started.....")
	dp.kafkaClient.Close()

	// Close Redis connection
	if err := dp.redisClient.Close(); err != nil {
		dp.log.Errorf("error closing redis connection: %v", err)
	}

	dp.log.Info("shut down completed.....")
}
