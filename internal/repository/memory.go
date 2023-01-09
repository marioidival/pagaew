package repository

import (
	"context"
	"fmt"
	"sync"
)

const (
	INVOICE_KEY     = "invoice_"
	LOG_INVOICE_KEY = "log_invoice_"
)

type LogMemory struct {
	store *sync.Map
}

type InvoiceMemory struct {
	store *sync.Map
}

func NewInvoiceMemoryRepository(store *sync.Map) InvoiceRepository {
	return &InvoiceMemory{
		store: store,
	}
}

func (m *InvoiceMemory) Get(ctx context.Context, debtID string) (*Invoice, error) {
	invoiceInterface, ok := m.store.Load(fmt.Sprintf("%s%s", INVOICE_KEY, debtID))
	if !ok {
		return nil, fmt.Errorf("%s debt not found", debtID)
	}
	invoice, ok := invoiceInterface.(Invoice)
	if !ok {
		return nil, fmt.Errorf("failed to load invoice: %s", debtID)
	}

	return &invoice, nil
}

func (m *InvoiceMemory) Save(ctx context.Context, invoices []Invoice) error {
	var possibleError error

	// emulate the behaivor of transaction
	possibleRemove := make([]Invoice, len(invoices))
	for _, invoice := range invoices {
		_, loaded := m.store.LoadOrStore(fmt.Sprintf("%s%s", INVOICE_KEY, invoice.DebtID), invoice)
		if loaded {
			possibleRemove = append(possibleRemove, invoice)
			possibleError = ErrInvoicesAlreadyExists
			break
		}
	}

	if possibleError != nil {
		for _, invoice := range possibleRemove {
			m.store.Delete(fmt.Sprintf("%s%s", INVOICE_KEY, invoice.DebtID))
		}
		return possibleError
	}

	for _, invoice := range invoices {
		m.store.Store(invoice.DebtID, invoice)
	}

	return nil
}

func (m *InvoiceMemory) Update(ctx context.Context, invoice Invoice) error {
	_, loaded := m.store.LoadOrStore(fmt.Sprintf("%s%s", INVOICE_KEY, invoice.DebtID), invoice)
	if !loaded {
		m.store.Delete(invoice.DebtID)
		return fmt.Errorf("%s invoice not exists", invoice.DebtID)
	}

	m.store.Store(invoice.DebtID, invoice)
	return nil
}

func NewLogMemoryRepository(store *sync.Map) LogRepository {
	return &LogMemory{
		store: store,
	}
}

func (l *LogMemory) Save(ctx context.Context, logInvoice LogInvoice) error {
	_, invoiceExists := l.store.Load(fmt.Sprintf("%s%s", INVOICE_KEY, logInvoice.DebtID))
	if !invoiceExists {
		return ErrInvoiceNotFound
	}

	_, loaded := l.store.LoadOrStore(fmt.Sprintf("%s%s", LOG_INVOICE_KEY, logInvoice.DebtID), logInvoice)
	if loaded {
		return ErrInvoiceAlreadyPaid
	}

	return nil
}
