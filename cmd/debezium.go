package cmd

import (
	"bytes"
	"encoding/json"
	"eventdrivensystem/configs"
	"eventdrivensystem/pkg/logger"
	"fmt"
	"io"
	"net/http"
)

// ConnectorManager handles Debezium connector operations
type ConnectorManager struct {
	config *configs.AppConfig
	log    logger.Logger
}

// ConnectorConfig holds Debezium connector settings
type ConnectorConfig struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// NewConnectorManager initializes a new ConnectorManager
func NewConnectorManager(config *configs.AppConfig, log logger.Logger) *ConnectorManager {
	return &ConnectorManager{
		config: config,
		log:    log,
	}
}

// Exists checks if the Debezium connector is already registered
func (cm *ConnectorManager) Exists(params ConnectorConfig) (bool, error) {
	resp, err := http.Get(cm.config.Debezium.KafkaConnectURL)
	if err != nil {
		return false, fmt.Errorf("failed to get connectors: %w", err)
	}
	defer resp.Body.Close()

	var connectors []string
	if err := json.NewDecoder(resp.Body).Decode(&connectors); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, conn := range connectors {
		if conn == params.Name {
			cm.log.Info("Connector already exists: ", params.Name)
			return true, nil
		}
	}
	return false, nil
}

// Register creates a new Debezium connector
func (cm *ConnectorManager) Register(params ConnectorConfig) error {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("failed to marshal connector data: %w", err)
	}

	resp, err := http.Post(cm.config.Debezium.KafkaConnectURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to register connector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to register connector, status: %d, response: %s", resp.StatusCode, string(body))
	}

	cm.log.Info("Connector registered successfully!")
	return nil
}

func (cm *ConnectorManager) Update(params ConnectorConfig) error {
	jsonData, err := json.Marshal(params.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal connector data: %w", err)
	}

	url := fmt.Sprintf("%s/%s/config", cm.config.Debezium.KafkaConnectURL, params.Name)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to update connector: %w", err)
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update connector: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update connector, status: %d, response: %s", resp.StatusCode, string(body))
	}

	cm.log.Info("Connector updated successfully!")
	return nil
}

// EnsureRegistration checks if the connector exists; if not, it registers it
func (cm *ConnectorManager) EnsureRegistration(params ConnectorConfig) error {
	exists, _ := cm.Exists(params)
	if exists {
		cm.log.Info("Connector already exists: ", params.Name)
		return nil
	}
	return cm.Register(params)
}

func (cm *ConnectorManager) RegisterOutboxConnector() error {
	dbzConnector := ConnectorConfig{
		Name: "outbox-connector",
		Config: map[string]interface{}{
			"connector.class":                        "io.debezium.connector.postgresql.PostgresConnector",
			"database.hostname":                      cm.config.Debezium.DBHost,
			"database.port":                          cm.config.Debezium.DBPort,
			"database.user":                          cm.config.Debezium.DBUsername,
			"database.password":                      cm.config.Debezium.DBPassword,
			"database.dbname":                        cm.config.Debezium.DBName,
			"database.server.name":                   "dbserver1",
			"table.include.list":                     cm.config.Debezium.TableName,
			"plugin.name":                            "pgoutput",
			"topic.prefix":                           "dbz",
			"key.converter":                          "org.apache.kafka.connect.json.JsonConverter",
			"value.converter":                        "org.apache.kafka.connect.json.JsonConverter",
			"transforms":                             "unwrap",
			"transforms.unwrap.type":                 "io.debezium.transforms.ExtractNewRecordState",
			"skipped.operations":                     "u,d",
			"snapshot.mode":                          "never",
			"event.processing.failure.handling.mode": "skip",
			"errors.tolerance":                       "all",
			"errors.deadletterqueue.topic.name":      "debezium-deadletter",
			"errors.deadletterqueue.context.headers.enable": "true",
			"errors.retry.delay.max.ms":                     "1000",
			"heartbeat.interval.ms":                         "10000",
		},
	}

	return cm.EnsureRegistration(dbzConnector)
}
