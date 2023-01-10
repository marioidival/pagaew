package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"

	"github.com/marioidival/pagaew/pkg/database"
)

type InvoiceMySQL struct {
	store *database.Client
}

type LogMySQL struct {
	store *database.Client
}

func NewInvoiceMySQLRepository(dbc *database.Client) InvoiceRepository {
	return &InvoiceMySQL{
		store: dbc,
	}
}

func (i *InvoiceMySQL) Get(ctx context.Context, debtID string) (*Invoice, error) {
	var invoice Invoice
	if err := i.store.QueryRow(ctx, "SELECT * FROM invoices WHERE debt_id = ?", debtID).Scan(&invoice); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvoiceNotFound
		}
		return nil, err
	}

	return &invoice, nil
}

func (i *InvoiceMySQL) Save(ctx context.Context, invoices []Invoice) error {
	tx, err := i.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmtIns, err := tx.PrepareContext(ctx, "INSERT INTO invoices(debt_id, debt_amount, debt_due_date, email, government_id, name) VALUES (?, ?, ?, ?, ?, ?)")

	for _, invoice := range invoices {
		_, err := stmtIns.ExecContext(
			ctx,
			invoice.DebtID,
			invoice.DebtAmount,
			invoice.DebtDueDate,
			invoice.Email,
			invoice.GovernmentID,
			invoice.Name,
		)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				return ErrInvoicesAlreadyExists
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (i *InvoiceMySQL) Update(ctx context.Context, invoice Invoice) error {
	tx, err := i.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var previousInvoice string
	err = tx.QueryRowContext(ctx, "SELECT * FROM invoices WHERE debt_id = ?", invoice.DebtID).Scan(&previousInvoice)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvoiceNotFound
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func NewLogMySQLRepository(dbc *database.Client) LogRepository {
	return &LogMySQL{
		store: dbc,
	}
}

func (l *LogMySQL) Save(ctx context.Context, logInvoice LogInvoice) error {
	tx, err := l.store.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var debtID string
	err = tx.QueryRowContext(ctx, "SELECT debt_id FROM invoices WHERE debt_id = ?", logInvoice.DebtID).Scan(&debtID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvoiceNotFound
		}

		return err
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO log_invoice(debt_id, paid_amount, paid_at, paid_by, status) VALUES (?, ?, ?, ?, ?)",
		logInvoice.DebtID,
		logInvoice.PaidAmount,
		logInvoice.PaidAt.Format("2006-01-02 15:04:05"),
		logInvoice.PaidBy,
		logInvoice.Status,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrInvoiceAlreadyPaid
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
