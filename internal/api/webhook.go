package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/marioidival/pagaew/internal/repository"
)

// Webhook handler to receive a JSON request notifying that invoice has been paid
func (s *Server) Webhook(ctx echo.Context) error {
	request := new(repository.LogInvoiceRequest)
	if err := ctx.Bind(request); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}

	logInvoice, err := toLogInvoice(*request)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{"message": err.Error()})
	}

	if err := s.logRepo.Save(ctx.Request().Context(), *logInvoice); err != nil {
		if errors.Is(err, repository.ErrInvoiceAlreadyPaid) {
			return ctx.JSON(http.StatusConflict, echo.Map{"message": err.Error()})
		}
		if errors.Is(err, repository.ErrInvoiceNotFound) {
			return ctx.JSON(http.StatusNotFound, echo.Map{"message": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}

	return ctx.JSON(http.StatusAccepted, echo.Map{"message": fmt.Sprintf("the invoice %s was queued to be paid", request.DebtID)})
}

func toLogInvoice(request repository.LogInvoiceRequest) (*repository.LogInvoice, error) {
	t, err := time.Parse(layoutWitTime, request.PaidAt)
	if err != nil {
		return nil, err
	}

	return &repository.LogInvoice{
		DebtID:     request.DebtID,
		PaidAt:     t,
		PaidAmount: request.PaidAmount,
		PaidBy:     request.PaidBy,
		Status:     repository.PENDING,
	}, nil
}
