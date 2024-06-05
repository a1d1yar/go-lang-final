package store

import (
	"database/sql"
	"fmt"
	"go-lang-final/internal/models"
)

type PaymentStore struct {
	DB *sql.DB
}

func NewPaymentStore(dsn string) (*PaymentStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &PaymentStore{DB: db}, nil
}

func (s *PaymentStore) CreatePayment(payment models.Payment) error {
	query := `INSERT INTO payments (id, amount, currency) VALUES ($1, $2, $3)`
	_, err := s.DB.Exec(query, payment.ID, payment.Amount, payment.Currency)
	return err
}

func (s *PaymentStore) GetPayment(id int64) (*models.Payment, error) {
	query := `SELECT id, amount, currency FROM payments WHERE id = $1`
	row := s.DB.QueryRow(query, id)

	var payment models.Payment
	if err := row.Scan(&payment.ID, &payment.Amount, &payment.Currency); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payment not found")
		}
		return nil, err
	}

	return &payment, nil
}

func (s *PaymentStore) UpdatePayment(id int64, payment models.Payment) error {
	query := `UPDATE payments SET amount = $2, currency = $3 WHERE id = $1`
	_, err := s.DB.Exec(query, id, payment.Amount, payment.Currency)
	return err
}

func (s *PaymentStore) DeletePayment(id int64) error {
	query := `DELETE FROM payments WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

func (s *PaymentStore) ListPayments(currency string, amount string, page int, pageSize int) ([]models.Payment, error) {
	query := `SELECT id, amount, currency FROM payments WHERE currency = $1 AND amount = $2 LIMIT $3 OFFSET $4`
	rows, err := s.DB.Query(query, currency, amount, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.ID, &payment.Amount, &payment.Currency); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
