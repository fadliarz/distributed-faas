package rpc

import (
	"context"
	"errors"

	"github.com/fadliarz/distributed-faas/services/billing-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/billing-service/domain/domain-core"
	billing_service_v1 "github.com/fadliarz/distributed-faas/services/billing-service/gen/go/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BillingServiceServer struct {
	billing_service_v1.UnimplementedBillingServiceServer
	handler *application.CommandHandler
}

func NewBillingServiceServer(handler *application.CommandHandler) *BillingServiceServer {
	return &BillingServiceServer{
		handler: handler,
	}
}

func (s *BillingServiceServer) Register(server *grpc.Server) {
	billing_service_v1.RegisterBillingServiceServer(server, s)
}

func (s *BillingServiceServer) GetBilling(ctx context.Context, req *billing_service_v1.GetBillingRequest) (*billing_service_v1.GetBillingResponse, error) {
	log.Info().Str("user_id", req.UserId).Msg("Received GetBilling request")

	query := &application.GetBillingQuery{
		UserID: req.UserId,
	}

	billing, err := s.handler.GetBilling(ctx, query)
	if err != nil {
		log.Error().Err(err).Str("user_id", req.UserId).Msg("Failed to get billing")

		if errors.Is(err, domain.ErrBillingNotFound) {
			return nil, status.Errorf(codes.NotFound, "billing not found for user ID: %s", req.UserId)
		}

		return nil, status.Errorf(codes.Internal, "failed to get billing: %v", err)
	}

	response := &billing_service_v1.GetBillingResponse{
		BillingId:  billing.BillingID.String(),
		UserId:     billing.UserID.String(),
		LastBilled: billing.LastBilled.Int64(),
		Amount:     billing.Amount.Int64(),
	}

	return response, nil
}
