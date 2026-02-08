package handlers

import (
	"encoding/json"
	"net/http"
	"tugas-session-1/models"
	"tugas-session-1/service"
)

// TransactionHandler handles transaction HTTP requests
type TransactionHandler struct {
	service service.TransactionService
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

// Checkout handles POST /transactions/checkout
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid JSON payload")
		return
	}

	transaction, err := h.service.Checkout(&req)
	if err != nil {
		// Check if it's a validation error
		if _, ok := err.(*models.ValidationError); ok {
			respondError(w, http.StatusBadRequest, "Bad Request", err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, transaction)
}

// GetTransactionByID handles GET /transactions/{id}
func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Bad Request", "Invalid transaction ID")
		return
	}

	transaction, err := h.service.GetByID(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if transaction == nil {
		respondError(w, http.StatusNotFound, "Not Found", "Transaction not found")
		return
	}

	respondJSON(w, http.StatusOK, transaction)
}

// GetAllTransactions handles GET /transactions
func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.service.GetAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, transactions)
}

// GetReport handles GET /report?start_date=2026-01-01&end_date=2026-02-01
func (h *TransactionHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	if startDate == "" {
		respondError(w, http.StatusBadRequest, "Bad Request", "start_date parameter is required")
		return
	}
	if endDate == "" {
		respondError(w, http.StatusBadRequest, "Bad Request", "end_date parameter is required")
		return
	}

	report, err := h.service.GetReport(startDate, endDate)
	if err != nil {
		// Check if it's a validation error
		if _, ok := err.(*models.ValidationError); ok {
			respondError(w, http.StatusBadRequest, "Bad Request", err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, report)
}

// GetTodayReport handles GET /api/report/hari-ini
func (h *TransactionHandler) GetTodayReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetTodayReport()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, report)
}

