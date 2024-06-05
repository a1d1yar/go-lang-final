package handlers

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/a1d1yar/go-lang-final/internal/models"
	"github.com/a1d1yar/go-lang-final/internal/store"
)

func RegisterGRPCHandlers(grpcServer *grpc.Server, store *store.PaymentStore) {
	RegisterPaymentServiceServer(grpcServer, &PaymentService{store: store})
}

type PaymentService struct {
	UnimplementedPaymentServiceServer
	store *store.PaymentStore
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	payment := models.Payment{
		Amount:   req.GetAmount(),
		Currency: req.GetCurrency(),
	}
	if err := s.store.CreatePayment(payment); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}
	return &CreatePaymentResponse{Success: true}, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, req *GetPaymentRequest) (*GetPaymentResponse, error) {
	id := int(req.GetId())
	payment, err := s.store.GetPayment(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}
	return &GetPaymentResponse{
		Id:       int32(payment.ID),
		Amount:   payment.Amount,
		Currency: payment.Currency,
	}, nil
}

func (s *PaymentService) UpdatePayment(ctx context.Context, req *UpdatePaymentRequest) (*UpdatePaymentResponse, error) {
	id := int(req.GetId())
	payment := models.Payment{
		Amount:   req.GetAmount(),
		Currency: req.GetCurrency(),
	}
	if err := s.store.UpdatePayment(id, payment); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update payment: %v", err)
	}
	return &UpdatePaymentResponse{Success: true}, nil
}

func (s *PaymentService) DeletePayment(ctx context.Context, req *DeletePaymentRequest) (*DeletePaymentResponse, error) {
	id := int(req.GetId())
	if err := s.store.DeletePayment(id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete payment: %v", err)
	}
	return &DeletePaymentResponse{Success: true}, nil
}
