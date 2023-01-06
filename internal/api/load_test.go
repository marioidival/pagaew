package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/marioidival/pagaew/internal/repository"
)

func loadTestFile(t *testing.T, filename string) (*bytes.Buffer, error) {
	t.Helper()

	file, err := os.Open(fmt.Sprintf("testdata/%s", filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, file)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func newServerInstance(t *testing.T) *Server {
	t.Helper()

	invoiceRepo := repository.NewInvoiceMemoryRepository()
	logRepo := repository.NewLogMemoryRepository()

	return NewServer(invoiceRepo, logRepo)
}

func TestLoad(t *testing.T) {
	t.Run("test basic load", func(t *testing.T) {
		e := echo.New()
		content, err := loadTestFile(t, "basicload.csv")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/load", content)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		server := newServerInstance(t)

		if assert.NoError(t, server.Load(c)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
	})

	t.Run("test wrong header load", func(t *testing.T) {
		e := echo.New()
		content, err := loadTestFile(t, "wrongheaderload.csv")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/load", content)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		server := newServerInstance(t)

		if assert.NoError(t, server.Load(c)) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("test wrong values in the file", func(t *testing.T) {
		e := echo.New()
		content, err := loadTestFile(t, "wrongvaluesload.csv")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/load", content)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		server := newServerInstance(t)

		if assert.NoError(t, server.Load(c)) {
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		}
	})

	t.Run("try save duplicate invoice", func(t *testing.T) {
		e := echo.New()
		content, err := loadTestFile(t, "duplicatedinvoices.csv")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/load", content)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		server := newServerInstance(t)

		if assert.NoError(t, server.Load(c)) {
			assert.Equal(t, http.StatusConflict, rec.Code)
		}
	})
}
