package repository

import (
	"database/sql"
	"fmt"
	"time"
	"tugas-session-1/models"
)

// TransactionRepository handles transaction data operations
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateCheckout creates a new transaction with details
func (r *TransactionRepository) CreateCheckout(totalAmount int, items []models.CheckoutItem, productPrices map[int]int) (*models.Transaction, error) {
	// Start a database transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert transaction
	var transactionID int
	var createdAt time.Time
	insertTxSQL := `
		INSERT INTO transactions (total_amount)
		VALUES ($1)
		RETURNING id, created_at
	`
	err = tx.QueryRow(insertTxSQL, totalAmount).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("error inserting transaction: %w", err)
	}

	// Insert transaction details and update product stock
	details := []models.TransactionDetail{}
	for _, item := range items {
		price := productPrices[item.ProductID]
		subtotal := price * item.Quantity

		// Insert detail
		var detailID int
		insertDetailSQL := `
			INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		err = tx.QueryRow(insertDetailSQL, transactionID, item.ProductID, item.Quantity, subtotal).Scan(&detailID)
		if err != nil {
			return nil, fmt.Errorf("error inserting transaction detail: %w", err)
		}

		// Update product stock
		updateStockSQL := `
			UPDATE products
			SET stock = stock - $1
			WHERE id = $2
		`
		_, err = tx.Exec(updateStockSQL, item.Quantity, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("error updating product stock: %w", err)
		}

		details = append(details, models.TransactionDetail{
			ID:            detailID,
			TransactionID: transactionID,
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			Subtotal:      subtotal,
		})
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	transaction := &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt,
		Details:     details,
	}

	return transaction, nil
}

// GetByID retrieves a transaction by ID with all details
func (r *TransactionRepository) GetByID(id int) (*models.Transaction, error) {
	// Get transaction header
	var transaction models.Transaction
	txQuery := `SELECT id, total_amount, created_at FROM transactions WHERE id = $1`
	err := r.db.QueryRow(txQuery, id).Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying transaction: %w", err)
	}

	// Get transaction details with product names
	detailsQuery := `
		SELECT td.id, td.transaction_id, td.product_id, p.name, td.quantity, td.subtotal
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		WHERE td.transaction_id = $1
		ORDER BY td.id
	`
	rows, err := r.db.Query(detailsQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error querying transaction details: %w", err)
	}
	defer rows.Close()

	var details []models.TransactionDetail
	for rows.Next() {
		var detail models.TransactionDetail
		err := rows.Scan(
			&detail.ID,
			&detail.TransactionID,
			&detail.ProductID,
			&detail.ProductName,
			&detail.Quantity,
			&detail.Subtotal,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning transaction detail: %w", err)
		}
		details = append(details, detail)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction details: %w", err)
	}

	transaction.Details = details
	return &transaction, nil
}

// GetAll retrieves all transactions with their details
func (r *TransactionRepository) GetAll() ([]models.Transaction, error) {
	// Get all transactions
	txQuery := `SELECT id, total_amount, created_at FROM transactions ORDER BY id DESC`
	rows, err := r.db.Query(txQuery)
	if err != nil {
		return nil, fmt.Errorf("error querying transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		err := rows.Scan(&tx.ID, &tx.TotalAmount, &tx.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	// Get details for each transaction
	for i := range transactions {
		detailsQuery := `
			SELECT td.id, td.transaction_id, td.product_id, p.name, td.quantity, td.subtotal
			FROM transaction_details td
			JOIN products p ON td.product_id = p.id
			WHERE td.transaction_id = $1
			ORDER BY td.id
		`
		detailRows, err := r.db.Query(detailsQuery, transactions[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error querying transaction details: %w", err)
		}

		var details []models.TransactionDetail
		for detailRows.Next() {
			var detail models.TransactionDetail
			err := detailRows.Scan(
				&detail.ID,
				&detail.TransactionID,
				&detail.ProductID,
				&detail.ProductName,
				&detail.Quantity,
				&detail.Subtotal,
			)
			if err != nil {
				detailRows.Close()
				return nil, fmt.Errorf("error scanning transaction detail: %w", err)
			}
			details = append(details, detail)
		}
		detailRows.Close()

		transactions[i].Details = details
	}

	return transactions, nil
}

// GetByDateRange retrieves transactions within a date range
func (r *TransactionRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
	// Get transactions in date range
	txQuery := `
		SELECT id, total_amount, created_at 
		FROM transactions 
		WHERE created_at >= $1 AND created_at < $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(txQuery, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error querying transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		err := rows.Scan(&tx.ID, &tx.TotalAmount, &tx.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning transaction: %w", err)
		}
		transactions = append(transactions, tx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	// Get details for each transaction
	for i := range transactions {
		detailsQuery := `
			SELECT td.id, td.transaction_id, td.product_id, p.name, td.quantity, td.subtotal
			FROM transaction_details td
			JOIN products p ON td.product_id = p.id
			WHERE td.transaction_id = $1
			ORDER BY td.id
		`
		detailRows, err := r.db.Query(detailsQuery, transactions[i].ID)
		if err != nil {
			return nil, fmt.Errorf("error querying transaction details: %w", err)
		}

		var details []models.TransactionDetail
		for detailRows.Next() {
			var detail models.TransactionDetail
			err := detailRows.Scan(
				&detail.ID,
				&detail.TransactionID,
				&detail.ProductID,
				&detail.ProductName,
				&detail.Quantity,
				&detail.Subtotal,
			)
			if err != nil {
				detailRows.Close()
				return nil, fmt.Errorf("error scanning transaction detail: %w", err)
			}
			details = append(details, detail)
		}
		detailRows.Close()

		transactions[i].Details = details
	}

	return transactions, nil
}
