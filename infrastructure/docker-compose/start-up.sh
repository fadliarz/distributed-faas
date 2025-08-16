#!/bin/bash

#
#
# ZOOKEEPER
#
#

echo "Starting Zookeeper"

docker-compose -f ./composes/common.yml -f ./composes/zookeeper.yml up -d --remove-orphans

sudo chown -R 1000:1000 ./composes/volumes

zookeeperCheckResult=$(echo ruok | nc localhost 2181)

while [[ ! $zookeeperCheckResult == "imok" ]]; do
    >&2 echo "Zookeeper is not running yet"
    sleep 2
    zookeeperCheckResult=$(echo ruok | nc localhost 2181)
done

#
#
# KAFKA CLUSTER
#
#

echo "Starting Kafka cluster"

docker-compose -f ./composes/common.yml -f ./composes/kafka_cluster.yml --profile "*" up -d

sudo chown -R 1000:1000 ./composes/volumes

kafkaCheckResult=$(kafkacat -L -b localhost:19092 | grep '2 brokers:')

while [[ ! $kafkaCheckResult == " 2 brokers:" ]]; do
    >&2 echo "Kafka cluster is not running yet"
    sleep 2
    kafkaCheckResult=$(kafkacat -L -b localhost:19092 | grep '2 brokers:')
done

#
#
# MONGO
#
#

echo "Starting Mongo"

docker-compose -f ./composes/common.yml -f ./composes/mongo.yml up -d

mongoCheckResult=$(docker exec -i distributed-faas-mongo mongosh --eval "db.adminCommand('ping')")

while [[ ! "$mongoCheckResult" == *'ok: 1'* ]]; do
    >&2 echo "Mongo is not running yet"
    sleep 2
    mongoCheckResult=$(docker exec -i distributed-faas-mongo mongosh --eval "db.adminCommand('ping')")
done

#
#
# DEBEZIUM CONNECTOR
#
#

echo "Creating Debezium connector"

# Wait for MongoDB to be fully up and running
sleep 10    

# invocation-db connector
curl --location --request POST 'localhost:8083/connectors' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "invocation-cdc",
    "config": {
        "connector.class": "io.debezium.connector.mongodb.MongoDbConnector",
        "mongodb.connection.string": "mongodb://admin:password@distributed-faas-mongo:27017/?replicaSet=rs0&directConnection=true",
        "topic.prefix": "cdc",
        "database.include.list": "invocation-db",
        "collection.include.list": "invocation-db.invocation",

        "key.converter": "org.apache.kafka.connect.json.JsonConverter",
        "key.converter.schemas.enable": false,
        "value.converter": "org.apache.kafka.connect.json.JsonConverter",
        "value.converter.schemas.enable": false,

        "transforms": "filter,unwrap",
        
        "transforms.filter.type": "io.debezium.transforms.Filter",
        "transforms.filter.language": "jsr223.groovy",
        "transforms.filter.condition": "value.op == 'c'",

        "transforms.unwrap.type": "io.debezium.connector.mongodb.transforms.ExtractNewDocumentState"
    }
}'

connectorCheckResult=$(curl --location --request GET 'localhost:8083/connectors/invocation-cdc/status')

while [[ ! "$connectorCheckResult" == *'"state":"RUNNING"'* ]]; do
    >&2 echo "Connector (invocation-cdc) is not running yet, waiting for it to start, connector check result: $connectorCheckResult"
    sleep 2
    connectorCheckResult=$(curl --location --request GET 'localhost:8083/connectors/invocation-cdc/status')
done

# checkpoint-db connector
curl --location --request POST 'localhost:8083/connectors' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "checkpoint-cdc",
    "config": {
        "connector.class": "io.debezium.connector.mongodb.MongoDbConnector",
        "mongodb.connection.string": "mongodb://admin:password@distributed-faas-mongo:27017/?replicaSet=rs0&directConnection=true",
        "topic.prefix": "cdc",
        "database.include.list": "checkpoint-db",
        "collection.include.list": "checkpoint-db.checkpoint",

        "key.converter": "org.apache.kafka.connect.json.JsonConverter",
        "key.converter.schemas.enable": false,
        "value.converter": "org.apache.kafka.connect.json.JsonConverter",
        "value.converter.schemas.enable": false,
			
        "transforms": "unwrap,filter,route",

        "transforms.unwrap.type": "io.debezium.connector.mongodb.transforms.ExtractNewDocumentState",

        "transforms.filter.type": "io.debezium.transforms.Filter",
        "transforms.filter.language": "jsr223.groovy",
        "transforms.filter.condition": "value.status == 'RETRYING' || value.status == 'SUCCESS'",

        "transforms.route.type": "io.debezium.transforms.ContentBasedRouter",
        "transforms.route.language": "jsr223.groovy",
        "transforms.route.topic.expression": "value.status == 'RETRYING' ? 'cdc.invocation-db.invocation' : null"
    }
}'

connectorCheckResult=$(curl --location --request GET 'localhost:8083/connectors/checkpoint-cdc/status')

while [[ ! "$connectorCheckResult" == *'"state":"RUNNING"'* ]]; do
    >&2 echo "Connector (checkpoint-cdc) is not running yet, waiting for it to start, connector check result: $connectorCheckResult"
    sleep 2
    connectorCheckResult=$(curl --location --request GET 'localhost:8083/connectors/checkpoint-cdc/status')
done

# cron-db connector
curl --location --request POST 'localhost:8083/connectors' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "cron-cdc",
    "config": {
            "connector.class": "io.debezium.connector.mongodb.MongoDbConnector",
            "mongodb.connection.string": "mongodb://admin:password@distributed-faas-mongo:27017/?replicaSet=rs0&directConnection=true",
            "topic.prefix": "cdc",
            "database.include.list": "cron-db",
            "collection.include.list": "cron-db.cron",

            "key.converter": "org.apache.kafka.connect.json.JsonConverter",
            "key.converter.schemas.enable": false,
            "value.converter": "org.apache.kafka.connect.json.JsonConverter",
            "value.converter.schemas.enable": false,

            "transforms": "filter,unwrap",

            "transforms.filter.type": "io.debezium.transforms.Filter",
            "transforms.filter.language": "jsr223.groovy",
            "transforms.filter.condition": "value.op == 'c' || value.op == 'u'",

            "transforms.unwrap.type": "io.debezium.connector.mongodb.transforms.ExtractNewDocumentState"
    }
}'

connectorCheckResult=$(curl --location --request GET 'localhost:8083/connectors/cron-cdc/status')

while [[ ! "$connectorCheckResult" == *'"state":"RUNNING"'* ]]; do
    >&2 echo "Connector (cron-cdc) is not running yet, waiting for it to start, connector check result: $connectorCheckResult"
    sleep 2
    connectorCheckResult=$(curl --location --request GET 'localhost:8083/connectors/cron-cdc/status')
done

#
#
# TEST DATA
#
#

echo "Inserting test data into MongoDB"

docker exec -it distributed-faas-mongo mongosh \
  --username admin \
  --password password \
  --authenticationDatabase admin \
  --eval "db.getSiblingDB('function-db').getCollection('function').insertOne({_id: '686da2aa129058624c7fb694', user_id: 'user-id-123', source_code_url: 'user-id-123/main.js'})"

docker exec -it distributed-faas-mongo mongosh \
  --username admin \
  --password password \
  --authenticationDatabase admin \
  --eval "db.getSiblingDB('machine-db').getCollection('machine').insertOne({_id: '686da3066a6680552bcddafb', address: 'distributed-faas-machine:50055', status: 'Available'})"

#
#
# SERVICES
#
#

# echo "Starting Services"

# docker-compose -f ./composes/common.yml -f ./composes/services.yml --profile "*" up --build -d