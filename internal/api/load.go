package api

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"github.com/marioidival/pagaew/internal/repository"
)

var expectedHeader = []string{"name", "governmentId", "email", "debtAmount", "debtDueDate", "debtId"}

// Load handler to receive a CVS document to save invoices to be paid
func (s *Server) Load(ctx echo.Context) error {
	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}

	r := csv.NewReader(bytes.NewReader(body))

	records, err := r.ReadAll()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": err.Error(),
		})
	}

	if !reflect.DeepEqual(expectedHeader, records[0]) {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid CSV header",
		})
	}

	invoices := make([]repository.Invoice, 0)

	for _, row := range records[1:] {
		if err := Validate(row); err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"message": err.Error()})
		}

		invoice, err := toInvoice(row)
		if err != nil {
			return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"message": err.Error()})
		}

		invoices = append(invoices, *invoice)
	}

	if err := s.invoiceRepo.Save(ctx.Request().Context(), invoices); err != nil {
		if errors.Is(err, repository.ErrInvoicesAlreadyExists) {
			return ctx.JSON(http.StatusConflict, echo.Map{"message": err.Error()})
		}

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, echo.Map{"message": "invoices saved successfully"})
}

func toInvoice(row []string) (*repository.Invoice, error) {
	debtAmountValue, err := toDebtAmount(row[debtAmount])
	if err != nil {
		return nil, err
	}
	debtDueDateValue, err := toDebtDueDate(row[debtDueDate])
	if err != nil {
		return nil, err
	}

	return &repository.Invoice{
		Name:         row[name],
		GovernmentID: row[governmentID],
		Email:        row[email],
		DebtAmount:   debtAmountValue,
		DebtDueDate:  debtDueDateValue,
		DebtID:       row[debtID],
	}, nil
}
