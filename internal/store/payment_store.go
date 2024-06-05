package store

import (
	"context"
	"fmt"
	"go-lang-final/internal/models"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PaymentStore struct {
	pool *pgxpool.Pool
}

func NewPaymentStore(dsn string) (*PaymentStore, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &PaymentStore{pool: pool}, nil
}

func (s *PaymentStore) CreatePayment(payment models.Payment) error {
	_, err := s.pool.Exec(context.Background(),
		"INSERT INTO payments (amount, currency) VALUES ($1, $2)",
		payment.Amount, payment.Currency)
	return err
}

func (s *PaymentStore) GetPayment(id int) (models.Payment, error) {
	var payment models.Payment
	err := s.pool.QueryRow(context.Background(),
		"SELECT id, amount, currency FROM payments WHERE id = $1", id).
		Scan(&payment.ID, &payment.Amount, &payment.Currency)
	return payment, err
}

func (s *PaymentStore) UpdatePayment(id int, payment models.Payment) error {
	_, err := s.pool.Exec(context.Background(),
		"UPDATE payments SET amount = $1, currency = $2 WHERE id = $3",
		payment.Amount, payment.Currency, id)
	return err
}

func (s *PaymentStore) DeletePayment(id int) error {
	_, err := s.pool.Exec(context.Background(), "DELETE FROM payments WHERE id = $1", id)
	return err
}

func (s *PaymentStore) ListPayments(filter string, sort string, page int, pageSize int) ([]models.Payment, error) {
	var payments []models.Payment
	var conditions []string
	var params []interface{}
	var sorting string

	if filter != "" {
		conditions = append(conditions, "currency ILIKE $"+strconv.Itoa(len(params)+1))
		params = append(params, "%"+filter+"%")
	}

	if sort != "" {
		sorting = fmt.Sprintf("ORDER BY %s", sort)
	} else {
		sorting = "ORDER BY id"
	}

	query := fmt.Sprintf(
		"SELECT id, amount, currency FROM payments WHERE %s %s LIMIT $%d OFFSET $%d",
		strings.Join(conditions, " AND "),
		sorting,
		len(params)+1,
		len(params)+2,
	)
	params = append(params, pageSize, (page-1)*pageSize)

	rows, err := s.pool.Query(context.Background(), query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment models.Payment
		if err := rows.Scan(&payment.ID, &payment.Amount, &payment.Currency); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}
