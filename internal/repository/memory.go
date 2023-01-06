package repository

import (
	"context"
	"fmt"
	"sync"
)

type LogMemory struct {
	store sync.Map
}

type InvoiceMemory struct {
	store sync.Map
}

func NewInvoiceMemoryRepository() InvoiceRepository {
	var m sync.Map
	return &InvoiceMemory{
		store: m,
	}
}

func (m *InvoiceMemory) Get(ctx context.Context, debtID uint) (*Invoice, error) {
	invoiceInterface, ok := m.store.Load(debtID)
	if !ok {
		return nil, fmt.Errorf("%d debt not found", debtID)
	}
	invoice, ok := invoiceInterface.(Invoice)
	if !ok {
		return nil, fmt.Errorf("failed to load invoice: %d", debtID)
	}

	return &invoice, nil
}

func (m *InvoiceMemory) Save(ctx context.Context, invoices []Invoice) error {
	var possibleError error

	// emulate the behaivor of transaction
	possibleRemove := make([]Invoice, len(invoices))
	for _, invoice := range invoices {
		_, loaded := m.store.LoadOrStore(invoice.DebtID, invoice)
		if loaded {
			possibleRemove = append(possibleRemove, invoice)
			possibleError = ErrInvoicesAlreadyExists
			break
		}
	}

	if possibleError != nil {
		for _, invoice := range possibleRemove {
			m.store.Delete(invoice.DebtID)
		}
		return possibleError
	}

	for _, invoice := range invoices {
		m.store.Store(invoice.DebtID, invoice)
	}

	return nil
}

func (m *InvoiceMemory) Update(ctx context.Context, invoice Invoice) error {
	_, loaded := m.store.LoadOrStore(invoice.DebtID, invoice)
	if !loaded {
		m.store.Delete(invoice.DebtID)
		return fmt.Errorf("%s invoice not exists", invoice.DebtID)
	}

	m.store.Store(invoice.DebtID, invoice)
	return nil
}

func NewLogMemoryRepository() LogRepository {
	var m sync.Map
	return &LogMemory{
		store: m,
	}
}

func (l *LogMemory) Save(ctx context.Context, logInvoice LogInvoice) error {
	_, loaded := l.store.LoadOrStore(logInvoice.DebtID, logInvoice)
	if loaded {
		return ErrInvoiceAlreadyPaid
	}

	return nil
}
