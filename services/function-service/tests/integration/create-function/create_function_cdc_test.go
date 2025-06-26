package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func (suite *CreateFunctionIntegrationTestSuite) setupDebeziumConnector() {
	connectPort, err := suite.debeziumConnectContainer.MappedPort(suite.ctx, "8083")
	if err != nil {
		suite.T().Logf("Failed to get Debezium Connect port: %v", err)
		return
	}

	connectHost, err := suite.debeziumConnectContainer.Host(suite.ctx)
	if err != nil {
		suite.T().Logf("Failed to get Debezium Connect host: %v", err)
		return
	}

	connectURL := fmt.Sprintf("http://%s:%d", connectHost, connectPort.Int())
	suite.T().Logf("Debezium Connect URL: %s", connectURL)

	suite.T().Logf("MongoURI: %s", os.Getenv("MONGO_URI"))

	// Create connector JSON configuration
	connectorName := fmt.Sprintf("mongodb-connector-test-fadli")
	// collectionIncludeList := fmt.Sprintf("%s\\.%s", suite.mongoDBName, suite.collectionName)

	// *** FIX: Use a double backslash to properly escape the dot for the JSON payload. ***
	// The regex requires a single backslash (\.), but JSON requires that backslash to be escaped (\\.).
	collectionIncludeList := fmt.Sprintf("%s\\\\.%s", suite.mongoDBName, suite.collectionName)

	log.Printf("mongodb.connection.string: mongodb://admin:password@%s:%d/?directConnection=true", suite.mongoHost, suite.mongoPort)

	configJSON := fmt.Sprintf(`{
		"name": "%s",
		"config": {
			"connector.class": "io.debezium.connector.mongodb.MongoDbConnector",
			"mongodb.connection.string": "%s",
			"topic.prefix": "test",
			"database.include.list": "%s",
			"collection.include.list": "%s",
			"key.converter": "org.apache.kafka.connect.json.JsonConverter",
			"value.converter": "org.apache.kafka.connect.json.JsonConverter",
			"key.converter.schemas.enable": false,
			"value.converter.schemas.enable": false,
			"snapshot.mode": "initial",
			"capture.mode": "change_streams_update_full"
		}
	}`, connectorName, fmt.Sprintf("mongodb://admin:password@%s:%d/?directConnection=true&replicaSet=rs0", suite.mongoHost, suite.mongoPort), suite.mongoDBName, collectionIncludeList)

	suite.T().Logf("mongodb.connection.uri: %s", fmt.Sprintf("mongodb://admin:password@%s:27017/", suite.mongoHost))

	// Create the HTTP request
	endpoint := fmt.Sprintf("%s/connectors", connectURL)
	suite.T().Logf("Registering Debezium connector at: %s", endpoint)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(configJSON))
	if err != nil {
		suite.T().Logf("Failed to create HTTP request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Implement retry logic
	maxRetries := 10
	backoffDelay := 2 * time.Second
	var resp *http.Response

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", endpoint, strings.NewReader(configJSON))

		if err != nil {
			suite.T().Logf("Failed to create HTTP request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		// Execute the request
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		resp, err = client.Do(req)

		if err != nil || resp.StatusCode >= 400 {
			if attempt < maxRetries {
				suite.T().Logf("Attempt %d: Failed to send HTTP request: %v", attempt, err)
				suite.T().Logf("Retrying in %v...", backoffDelay)
				time.Sleep(backoffDelay)

				if err == nil {
					body, _ := io.ReadAll(resp.Body)
					suite.T().Logf("Response status: %d, body: %s", resp.StatusCode, string(body))
				}

				backoffDelay *= 2
				continue
			}
			return
		}
		break
	}

	if resp == nil {
		suite.T().Logf("Failed to connect to Debezium after %d attempts", maxRetries)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		suite.T().Logf("Successfully registered Debezium connector: %s", connectorName)
	} else {
		body, _ := io.ReadAll(resp.Body)
		suite.T().Fatalf("Failed to register Debezium connector. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	// Wait for the connector to be ready
	suite.waitForDebeziumConnectorReady(connectURL, connectorName)
}

func (suite *CreateFunctionIntegrationTestSuite) waitForDebeziumConnectorReady(connectURL, connectorName string) {
	maxWaitTime := 30 * time.Second
	checkInterval := 2 * time.Second
	deadline := time.Now().Add(maxWaitTime)

	suite.T().Logf("Waiting for Debezium connector '%s' to be ready...", connectorName)

	for time.Now().Before(deadline) {
		// Check connector status
		statusURL := fmt.Sprintf("%s/connectors/%s/status", connectURL, connectorName)
		resp, err := http.Get(statusURL)
		if err != nil {
			suite.T().Logf("Error checking connector status: %v", err)
			time.Sleep(checkInterval)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			suite.T().Logf("Error reading status response: %v", err)
			time.Sleep(checkInterval)
			continue
		}

		statusStr := string(body)
		suite.T().Logf("Connector status: %s", statusStr)

		// Check if connector is running
		if strings.Contains(statusStr, `"state":"RUNNING"`) {
			suite.T().Logf("Debezium connector '%s' is now running!", connectorName)
			// Give it a bit more time to start streaming
			time.Sleep(5 * time.Second)
			return
		}

		time.Sleep(checkInterval)
	}

	suite.T().Fatalf("Debezium connector '%s' did not become ready within %v", connectorName, maxWaitTime)
}
