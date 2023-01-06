package repository

import (
	"context"
	"errors"
	"time"
)

var ErrInvoicesAlreadyExists = errors.New("invoices already exists")

type Invoice struct {
	Name         string
	GovernmentID string
	Email        string
	DebtAmount   float64
	DebtDueDate  time.Time
	DebtID       string
}

type InvoiceRepository interface {
	Get(ctx context.Context, debtID uint) (*Invoice, error)
	Save(ctx context.Context, invoices []Invoice) error
	Update(ctx context.Context, invoice Invoice) error
}

type LogInvoice struct {
	DebtID     uint      `json:"debtId"`
	PaidAt     time.Time `json:"paidAt"`
	PaidAmount float64   `json:"paidAmount"`
	PaidBy     string    `json:"paidBy"`
}

type LogRepository interface {
	Save(ctx context.Context, logInvoice LogInvoice) error
}
