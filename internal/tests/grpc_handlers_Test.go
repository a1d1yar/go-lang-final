package tests

import (
	"context"
	"go-lang-final/internal/store"
	"go-lang-final/proto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type PaymentService struct {
	proto.UnimplementedPaymentServiceServer
	store *store.PaymentStore
}

func TestGRPCCreatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO payments").
		WithArgs(1, 100.0, "USD").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := &store.PaymentStore{DB: db}

	h := &PaymentService{store: s}
	req := &proto.CreatePaymentRequest{Id: 1, Amount: 100.0, Currency: "USD"}

	res, err := h.CreatePayment(context.Background(), req)
	assert.NoError(t, err)
	assert.True(t, res.GetSuccess())
}

func TestGRPCGetPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "amount", "currency"}).
		AddRow(1, 100.0, "USD")

	mock.ExpectQuery("SELECT id, amount, currency FROM payments WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	s := &store.PaymentStore{DB: db}
	h := &PaymentService{store: s}
	req := &proto.GetPaymentRequest{Id: 1}

	res, err := h.GetPayment(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, &proto.GetPaymentResponse{Id: 1, Amount: 100.0, Currency: "USD"}, res)
}

func TestGRPCUpdatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("UPDATE payments SET amount = ?, currency = ? WHERE id = ?").
		WithArgs(100.0, "USD", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := &store.PaymentStore{DB: db}
	h := &PaymentService{store: s}
	req := &proto.UpdatePaymentRequest{Id: 1, Amount: 100.0, Currency: "USD"}

	res, err := h.UpdatePayment(context.Background(), req)
	assert.NoError(t, err)
	assert.True(t, res.GetSuccess())
}

func TestGRPCDeletePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("DELETE FROM payments WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := &store.PaymentStore{DB: db}
	h := &PaymentService{store: s}
	req := &proto.DeletePaymentRequest{Id: 1}

	res, err := h.DeletePayment(context.Background(), req)
	assert.NoError(t, err)
	assert.True(t, res.GetSuccess())
}

func TestGRPCListPayments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "amount", "currency"}).
		AddRow(1, 100.0, "USD").
		AddRow(2, 200.0, "USD")

	mock.ExpectQuery("SELECT id, amount, currency FROM payments WHERE currency = ? AND amount = ? LIMIT ? OFFSET ?").
		WithArgs("USD", "100.00", 10, 0).
		WillReturnRows(rows)

	s := &store.PaymentStore{DB: db}
	h := &PaymentService{store: s}
	req := &proto.ListPaymentsRequest{Currency: "USD", Amount: 100.0, Page: 1, PageSize: 10}

	res, err := h.ListPayments(context.Background(), req)
	assert.NoError(t, err)
	assert.Len(t, res.GetPayments(), 2)
	assert.Equal(t, &proto.Payment{Id: 1, Amount: 100.0, Currency: "USD"}, res.GetPayments()[0])
	assert.Equal(t, &proto.Payment{Id: 2, Amount: 200.0, Currency: "USD"}, res.GetPayments()[1])
}
