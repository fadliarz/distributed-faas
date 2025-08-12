module github.com/fadliarz/distributed-faas/services/billing-calculator-service

go 1.24.4

require (
	github.com/fadliarz/distributed-faas/common v0.0.0-00010101000000-000000000000
	github.com/fadliarz/distributed-faas/infrastructure/kafka v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.34.0
	go.mongodb.org/mongo-driver v1.17.1
)

replace github.com/fadliarz/distributed-faas/common => ../../common

replace github.com/fadliarz/distributed-faas/infrastructure/kafka => ../../infrastructure/kafka
