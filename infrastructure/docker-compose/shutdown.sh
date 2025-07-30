echo "Shutdown zookeeper"

docker-compose -f ./composes/common.yml -f ./composes/zookeeper.yml down

sleep 5

echo "Shutdown kafka cluster"

docker-compose -f ./composes/common.yml -f ./composes/kafka_cluster.yml  --profile "*" down

sleep 5

echo "Shutdown mongo"

docker-compose -f ./composes/common.yml -f ./composes/mongo.yml --profile "*" down

sleep 5

echo "Shutdown services"

docker-compose -f ./composes/common.yml -f ./composes/services.yml --profile "*" down

echo "Deleting Kafka and Zookeeper volumes"

yes | sudo rm -r ./composes/volumes/zookeeper

yes | sudo rm -r ./composes/volumes/kafka