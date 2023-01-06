package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	server := newServerInstance(t)

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
}
