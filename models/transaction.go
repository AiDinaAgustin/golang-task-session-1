package models

import "time"

// Transaction represents a transaction entity
type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

// TransactionDetail represents a transaction detail entity
type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

// CheckoutItem represents an item in the checkout request
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// CheckoutRequest represents the request body for creating a transaction
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

// ReportResponse represents the response for report endpoint
type ReportResponse struct {
	StartDate        string        `json:"start_date"`
	EndDate          string        `json:"end_date"`
	TotalTransactions int          `json:"total_transactions"`
	TotalRevenue     int           `json:"total_revenue"`
	Transactions     []Transaction `json:"transactions"`
}

// Validate validates the checkout request
func (r *CheckoutRequest) Validate() error {
	if len(r.Items) == 0 {
		return &ValidationError{Field: "items", Message: "items cannot be empty"}
	}
	
	for i, item := range r.Items {
		if item.ProductID <= 0 {
			return &ValidationError{Field: "items", Message: "product_id must be greater than 0"}
		}
		if item.Quantity <= 0 {
			return &ValidationError{Field: "items", Message: "quantity must be greater than 0"}
		}
		
		// Check for duplicate product IDs
		for j := i + 1; j < len(r.Items); j++ {
			if r.Items[j].ProductID == item.ProductID {
				return &ValidationError{Field: "items", Message: "duplicate product_id not allowed"}
			}
		}
	}
	
	return nil
}
