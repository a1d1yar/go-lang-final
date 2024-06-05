package handlers

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-lang-final/internal/models"
	"go-lang-final/internal/store"
	"go-lang-final/proto"
)

func RegisterGRPCHandlers(grpcServer *grpc.Server, store *store.PaymentStore) {
	proto.RegisterPaymentServiceServer(grpcServer, &PaymentService{store: store})
}

type PaymentService struct {
	proto.UnimplementedPaymentServiceServer
	store *store.PaymentStore
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *proto.CreatePaymentRequest) (*proto.CreatePaymentResponse, error) {
	payment := models.Payment{
		ID:       req.GetId(),
		Amount:   req.GetAmount(),
		Currency: req.GetCurrency(),
	}

	err := s.store.CreatePayment(payment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}

	return &proto.CreatePaymentResponse{Success: true}, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, req *proto.GetPaymentRequest) (*proto.GetPaymentResponse, error) {
	payment, err := s.store.GetPayment(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	return &proto.GetPaymentResponse{
		Id:       payment.ID,
		Amount:   payment.Amount,
		Currency: payment.Currency,
	}, nil
}

func (s *PaymentService) UpdatePayment(ctx context.Context, req *proto.UpdatePaymentRequest) (*proto.UpdatePaymentResponse, error) {
	payment := models.Payment{
		ID:       req.GetId(),
		Amount:   req.GetAmount(),
		Currency: req.GetCurrency(),
	}

	err := s.store.UpdatePayment(req.GetId(), payment)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update payment: %v", err)
	}

	return &proto.UpdatePaymentResponse{Success: true}, nil
}

func (s *PaymentService) DeletePayment(ctx context.Context, req *proto.DeletePaymentRequest) (*proto.DeletePaymentResponse, error) {
	err := s.store.DeletePayment(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete payment: %v", err)
	}

	return &proto.DeletePaymentResponse{Success: true}, nil
}

func (s *PaymentService) ListPayments(ctx context.Context, req *proto.ListPaymentsRequest) (*proto.ListPaymentsResponse, error) {
	payments, err := s.store.ListPayments(req.GetCurrency(), fmt.Sprintf("%.2f", req.GetAmount()), int(req.GetPage()), int(req.GetPageSize()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list payments: %v", err)
	}

	var paymentProtos []*proto.Payment
	for _, payment := range payments {
		paymentProtos = append(paymentProtos, &proto.Payment{
			Id:       payment.ID,
			Amount:   payment.Amount,
			Currency: payment.Currency,
		})
	}

	return &proto.ListPaymentsResponse{Payments: paymentProtos}, nil
}
