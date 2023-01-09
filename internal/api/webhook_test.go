package api

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/marioidival/pagaew/internal/repository"
)

func fixtureSaveInvoice(t *testing.T, server *Server) []repository.Invoice {
	t.Helper()

	newInvoices := make([]repository.Invoice, 0)

	newInvoice := repository.Invoice{
		Name:         "testing",
		GovernmentID: "12133333333",
		Email:        "myemail@gmail.com",
		DebtAmount:   1333.11,
		DebtDueDate:  time.Now(),
		DebtID:       "1234",
	}

	newInvoices = append(newInvoices, newInvoice)

	newInvoice = repository.Invoice{
		Name:         "John Doe",
		GovernmentID: "12345678900",
		Email:        "myemail@gmail.com",
		DebtAmount:   1333.11,
		DebtDueDate:  time.Now(),
		DebtID:       "8291",
	}
	newInvoices = append(newInvoices, newInvoice)

	if err := server.invoiceRepo.Save(context.Background(), newInvoices); err != nil {
		t.Fail()
	}

	return newInvoices
}

func TestWebhook(t *testing.T) {
	server := newServerInstance(t)

	_ = fixtureSaveInvoice(t, server)

	t.Run("test save weebhook", func(t *testing.T) {
		e := echo.New()

		body := bytes.NewBuffer([]byte(`{"debtId": "8291", "paidAt": "2022-06-09 10:00:00", "paidAmount": 100000.00, "paidBy": "John Doe"}`))

		req := httptest.NewRequest(http.MethodPost, "/webhook", body)
		req.Header.Add(echo.HeaderContentType, "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		if assert.NoError(t, server.Webhook(c)) {
			assert.Equal(t, http.StatusAccepted, rec.Code)
			assert.True(t, strings.Contains(rec.Body.String(), "8291"))
		}
	})

	t.Run("test save same weebhook should be failed", func(t *testing.T) {
		e := echo.New()

		body := bytes.NewBuffer([]byte(`{"debtId": "8291", "paidAt": "2022-06-09 10:00:00", "paidAmount": 100000.00, "paidBy": "John Doe"}`))

		req := httptest.NewRequest(http.MethodPost, "/webhook", body)
		req.Header.Add(echo.HeaderContentType, "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		if assert.NoError(t, server.Webhook(c)) {
			assert.Equal(t, http.StatusConflict, rec.Code)
		}
	})

	t.Run("test save invalid weebhook values", func(t *testing.T) {
		e := echo.New()

		body := bytes.NewBuffer([]byte(`{"debtId": "8291", "paidAt": "2022-06-0910:00:00", "paidAmount": "100000.00", "paidBy": "John Doe"}`))

		req := httptest.NewRequest(http.MethodPost, "/webhook", body)
		req.Header.Add(echo.HeaderContentType, "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		if assert.NoError(t, server.Webhook(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("test save a weebhook that not exists", func(t *testing.T) {
		e := echo.New()

		body := bytes.NewBuffer([]byte(`{"debtId": "1335", "paidAt": "2022-06-09 10:00:00", "paidAmount": 100000.00, "paidBy": "John Doe"}`))

		req := httptest.NewRequest(http.MethodPost, "/webhook", body)
		req.Header.Add(echo.HeaderContentType, "application/json")
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		if assert.NoError(t, server.Webhook(c)) {
			log.Println(rec.Body.String())
			assert.Equal(t, http.StatusNotFound, rec.Code)
		}
	})
}
