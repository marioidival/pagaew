package repository

import (
	"context"
	"errors"
	"time"
)

const (
	DONE    = "DONE"
	ERROR   = "ERROR"
	PENDING = "PENDING"
)

var (
	ErrInvoicesAlreadyExists = errors.New("invoices already exists")
	ErrInvoiceAlreadyPaid    = errors.New("invoice already paid")
)

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

type LogInvoiceRequest struct {
	DebtID     string  `json:"debtId"`
	PaidAt     string  `json:"paidAt"`
	PaidAmount float64 `json:"paidAmount"`
	PaidBy     string  `json:"paidBy"`
}

type LogInvoice struct {
	DebtID     string
	PaidAt     time.Time
	PaidAmount float64
	PaidBy     string
	Status     string
}

type LogRepository interface {
	Save(ctx context.Context, logInvoice LogInvoice) error
}
