package rpc

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/charge-service/domain/application-service"
	charge_service_v1 "github.com/fadliarz/distributed-faas/services/charge-service/gen/go/charge-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChargeServer struct {
	charge_service_v1.UnimplementedChargeServiceServer

	handler *application.CommandHandler
}

func NewChargeServer(handler *application.CommandHandler) *ChargeServer {
	return &ChargeServer{
		handler: handler,
	}
}

func (s *ChargeServer) Register(server *grpc.Server) {
	charge_service_v1.RegisterChargeServiceServer(server, s)
}

func (s *ChargeServer) CreateCharge(ctx context.Context, req *charge_service_v1.CreateChargeRequest) (*charge_service_v1.CreateChargeResponse, error) {
	cmd := &application.CreateChargeCommand{
		UserID:    req.UserId,
		ServiceID: req.ServiceId,
		Amount:    req.Amount,
	}

	if err := s.handler.CreateCharge(ctx, cmd); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create charge: %v", err)
	}

	return &charge_service_v1.CreateChargeResponse{
		Status:  "success",
		Message: "Charge created and queued for processing",
	}, nil
}
