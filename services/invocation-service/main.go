package main

import (
	"fmt"
	"log"
	"net"

	"github.com/fadliarz/services/invocation-service/application"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := ":50052"

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("gRPC server listening on %s\n", port)

	grpcServer := grpc.NewServer()

	functionServer := application.NewInvocationServer()
	functionServer.Register(grpcServer)

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
