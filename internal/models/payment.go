package models

type Payment struct {
	ID       int64   `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
