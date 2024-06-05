package tests

import (
	"go-lang-final/internal/models"
	"go-lang-final/internal/store"

	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("INSERT INTO payments").
		WithArgs(1, 100.0, "USD").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := &store.PaymentStore{db: db}
	err = s.CreatePayment(models.Payment{ID: 1, Amount: 100.0, Currency: "USD"})
	assert.NoError(t, err)
}

func TestGetPayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "amount", "currency"}).
		AddRow(1, 100.0, "USD")

	mock.ExpectQuery("SELECT id, amount, currency FROM payments WHERE id = ?").
		WithArgs(1).
		WillReturnRows(rows)

	s := &store.PaymentStore{db: db}
	payment, err := s.GetPayment(1)
	assert.NoError(t, err)
	assert.Equal(t, &models.Payment{ID: 1, Amount: 100.0, Currency: "USD"}, payment)
}

func TestUpdatePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("UPDATE payments SET amount = ?, currency = ? WHERE id = ?").
		WithArgs(100.0, "USD", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := &store.PaymentStore{db: db}
	err = s.UpdatePayment(1, models.Payment{ID: 1, Amount: 100.0, Currency: "USD"})
	assert.NoError(t, err)
}

func TestDeletePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec("DELETE FROM payments WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := &store.PaymentStore{db: db}
	err = s.DeletePayment(1)
	assert.NoError(t, err)
}

func TestListingPayments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "amount", "currency"}).
		AddRow(1, 100.0, "USD").
		AddRow(2, 200.0, "USD")

	mock.ExpectQuery("SELECT id, amount, currency FROM payments WHERE currency = ? AND amount = ? LIMIT ? OFFSET ?").
		WithArgs("USD", "100.00", 10, 0).
		WillReturnRows(rows)

	s := &store.PaymentStore{db: db}
	payments, err := s.ListPayments("USD", "100.00", 1, 10)
	assert.NoError(t, err)
	assert.Len(t, payments, 2)
	assert.Equal(t, []models.Payment{
		{ID: 1, Amount: 100.0, Currency: "USD"},
		{ID: 2, Amount: 200.0, Currency: "USD"},
	}, payments)
}
