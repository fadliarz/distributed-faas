package cdc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type DebeziumConnectorConfig struct {
	Name                  string
	ConnectorClass        string
	TasksMax              int
	MongodbConnectionURI  string
	MongodbUser           string
	MongodbPassword       string
	MongodbHosts          string
	MongodbName           string
	CollectionIncludeList string
	TopicPrefix           string
	KeyConverter          string
	ValueConverter        string
	SchemaRegistryURL     string
	SnapshotMode          string
	TransformsSMT         []string
	TransformConfigs      map[string]map[string]string
}

func NewMongoDBDebeziumConfig(name, mongoURI, dbName, collectionList, topicPrefix, schemaRegistryURL string) *DebeziumConnectorConfig {
	user, password, hosts := parseMongoURI(mongoURI)

	return &DebeziumConnectorConfig{
		Name:                  name,
		ConnectorClass:        "io.debezium.connector.mongodb.MongoDbConnector",
		TasksMax:              1,
		MongodbConnectionURI:  mongoURI,
		MongodbUser:           user,
		MongodbPassword:       password,
		MongodbHosts:          hosts,
		MongodbName:           dbName,
		CollectionIncludeList: collectionList,
		TopicPrefix:           topicPrefix,
		KeyConverter:          "io.confluent.connect.avro.AvroConverter",
		ValueConverter:        "io.confluent.connect.protobuf.ProtobufConverter",
		SchemaRegistryURL:     schemaRegistryURL,
		SnapshotMode:          "initial",
		TransformsSMT:         []string{"unwrap", "extractKey"},
		TransformConfigs: map[string]map[string]string{
			"unwrap": {
				"type":                 "io.debezium.transforms.ExtractNewRecordState",
				"delete.handling.mode": "rewrite",
				"add.fields":           "op,ts_ms",
			},
			"extractKey": {
				"type":  "org.apache.kafka.connect.transforms.ExtractField$Key",
				"field": "id",
			},
		},
	}
}

func (c *DebeziumConnectorConfig) ToJSON() (string, error) {
	config := map[string]interface{}{
		"name": c.Name,
		"config": map[string]interface{}{
			"connector.class":                     c.ConnectorClass,
			"tasks.max":                           c.TasksMax,
			"mongodb.connection.uri":              c.MongodbConnectionURI,
			"mongodb.user":                        c.MongodbUser,
			"mongodb.password":                    c.MongodbPassword,
			"mongodb.hosts":                       c.MongodbHosts,
			"mongodb.name":                        c.MongodbName,
			"collection.include.list":             c.CollectionIncludeList,
			"topic.prefix":                        c.TopicPrefix,
			"key.converter":                       c.KeyConverter,
			"value.converter":                     c.ValueConverter,
			"key.converter.schema.registry.url":   c.SchemaRegistryURL,
			"value.converter.schema.registry.url": c.SchemaRegistryURL,
			"snapshot.mode":                       c.SnapshotMode,
		},
	}

	if len(c.TransformsSMT) > 0 {
		config["config"].(map[string]interface{})["transforms"] = strings.Join(c.TransformsSMT, ",")
		for name, transformConfig := range c.TransformConfigs {
			for key, value := range transformConfig {
				config["config"].(map[string]interface{})[fmt.Sprintf("transforms.%s.%s", name, key)] = value
			}
		}
	}

	jsonData, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Debezium config to JSON: %w", err)
	}

	return string(jsonData), nil
}

func (c *DebeziumConnectorConfig) RegisterConnector(connectURL string) error {
	jsonConfig, err := c.ToJSON()
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("%s/connectors", connectURL)
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(jsonConfig))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to register connector, HTTP status: %d", resp.StatusCode)
	}

	log.Printf("Successfully registered Debezium connector: %s", c.Name)
	return nil
}

func parseMongoURI(uri string) (user, password, hosts string) {
	// mongodb://username:password@host1:port1,host2:port2/database

	// Extract user and password
	if strings.Contains(uri, "@") {
		credentials := strings.Split(strings.Split(uri, "@")[0], "//")[1]
		parts := strings.Split(credentials, ":")
		if len(parts) > 1 {
			user = parts[0]
			password = parts[1]
		} else {
			user = parts[0]
		}
	}

	// Extract hosts
	if strings.Contains(uri, "@") {
		hostsStr := strings.Split(strings.Split(uri, "@")[1], "/")[0]
		hosts = hostsStr
	} else {
		hostsStr := strings.Split(strings.Split(uri, "//")[1], "/")[0]
		hosts = hostsStr
	}

	return
}
