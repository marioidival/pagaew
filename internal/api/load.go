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

		invoices = append(invoices, toInvoice(row))
	}

	if err := s.invoiceRepo.Save(ctx.Request().Context(), invoices); err != nil {
		if errors.Is(err, repository.ErrInvoicesAlreadyExists) {
			return ctx.JSON(http.StatusConflict, echo.Map{"message": err.Error()})
		}

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, echo.Map{"message": "invoices saved successfully"})
}

func toInvoice(row []string) repository.Invoice {
	debtAmountValue, _ := toDebtAmount(row[debtAmount])
	debtDueDateValue, _ := toDebtDueDate(row[debtDueDate])

	return repository.Invoice{
		Name:         row[name],
		GovernmentID: row[governmentID],
		Email:        row[email],
		DebtAmount:   debtAmountValue,
		DebtDueDate:  debtDueDateValue,
		DebtID:       row[debtID],
	}
}
