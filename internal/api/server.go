package api

import "github.com/marioidival/pagaew/internal/repository"

type Server struct {
	invoiceRepo repository.InvoiceRepository
	logRepo repository.LogRepository
}

func NewServer(invoiceRepository repository.InvoiceRepository, logRepostiroy repository.LogRepository) *Server {
	return &Server{
		invoiceRepo: invoiceRepository,
		logRepo: logRepostiroy,
	}
}
