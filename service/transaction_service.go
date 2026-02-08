package service

import (
	"fmt"
	"time"
	"tugas-session-1/models"
	"tugas-session-1/repository"
)

// TransactionService defines the interface for transaction business logic
type TransactionService interface {
	Checkout(req *models.CheckoutRequest) (*models.Transaction, error)
	GetByID(id int) (*models.Transaction, error)
	GetAll() ([]models.Transaction, error)
	GetReport(startDate, endDate string) (*models.ReportResponse, error)
	GetTodayReport() (*models.ReportResponse, error)
}

// transactionService implements TransactionService interface
type transactionService struct {
	transactionRepo *repository.TransactionRepository
	productRepo     *repository.ProductRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(transactionRepo *repository.TransactionRepository, productRepo *repository.ProductRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

// Checkout processes a checkout request
func (s *transactionService) Checkout(req *models.CheckoutRequest) (*models.Transaction, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Validate all products exist and have sufficient stock
	totalAmount := 0
	productPrices := make(map[int]int)

	for _, item := range req.Items {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("error fetching product: %w", err)
		}
		if product == nil {
			return nil, &models.ValidationError{
				Field:   "items",
				Message: fmt.Sprintf("product with id %d not found", item.ProductID),
			}
		}

		// Check stock availability
		if product.Stock < item.Quantity {
			return nil, &models.ValidationError{
				Field:   "items",
				Message: fmt.Sprintf("insufficient stock for product %s (available: %d, requested: %d)", product.Name, product.Stock, item.Quantity),
			}
		}

		// Calculate subtotal using price as-is (already in rupiah)
		price := int(product.Price)
		productPrices[item.ProductID] = price
		totalAmount += price * item.Quantity
	}

	// Create transaction
	transaction, err := s.transactionRepo.CreateCheckout(totalAmount, req.Items, productPrices)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction: %w", err)
	}

	// Fetch the complete transaction with product names
	completeTransaction, err := s.transactionRepo.GetByID(transaction.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching complete transaction: %w", err)
	}

	return completeTransaction, nil
}

// GetByID retrieves a transaction by ID
func (s *transactionService) GetByID(id int) (*models.Transaction, error) {
	return s.transactionRepo.GetByID(id)
}

// GetAll retrieves all transactions
func (s *transactionService) GetAll() ([]models.Transaction, error) {
	transactions, err := s.transactionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Return empty array instead of nil
	if transactions == nil {
		return []models.Transaction{}, nil
	}

	return transactions, nil
}

// GetReport generates a report for transactions within a date range
func (s *transactionService) GetReport(startDate, endDate string) (*models.ReportResponse, error) {
	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, &models.ValidationError{
			Field:   "start_date",
			Message: "invalid date format, expected YYYY-MM-DD",
		}
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, &models.ValidationError{
			Field:   "end_date",
			Message: "invalid date format, expected YYYY-MM-DD",
		}
	}

	// Add one day to end date to include transactions on that day
	end = end.AddDate(0, 0, 1)

	// Validate date range
	if end.Before(start) {
		return nil, &models.ValidationError{
			Field:   "end_date",
			Message: "end_date must be after start_date",
		}
	}

	// Get transactions in date range
	transactions, err := s.transactionRepo.GetByDateRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("error fetching transactions: %w", err)
	}

	// Calculate summary
	totalRevenue := 0
	for _, tx := range transactions {
		totalRevenue += tx.TotalAmount
	}

	// Return empty array instead of nil
	if transactions == nil {
		transactions = []models.Transaction{}
	}

	report := &models.ReportResponse{
		StartDate:         startDate,
		EndDate:           endDate,
		TotalTransactions: len(transactions),
		TotalRevenue:      totalRevenue,
		Transactions:      transactions,
	}

	return report, nil
}

// GetTodayReport generates a report for today's transactions
func (s *transactionService) GetTodayReport() (*models.ReportResponse, error) {
	// Get today's date in local timezone
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	// Get transactions for today
	transactions, err := s.transactionRepo.GetByDateRange(startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("error fetching today's transactions: %w", err)
	}

	// Calculate summary
	totalRevenue := 0
	for _, tx := range transactions {
		totalRevenue += tx.TotalAmount
	}

	// Return empty array instead of nil
	if transactions == nil {
		transactions = []models.Transaction{}
	}

	todayDateStr := startOfDay.Format("2006-01-02")

	report := &models.ReportResponse{
		StartDate:         todayDateStr,
		EndDate:           todayDateStr,
		TotalTransactions: len(transactions),
		TotalRevenue:      totalRevenue,
		Transactions:      transactions,
	}

	return report, nil
}

