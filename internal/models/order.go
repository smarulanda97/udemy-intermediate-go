package models

import (
	"context"
	"time"
)

// Order is the type for all orders
type Order struct {
	ID            int       `json:"id"`
	CustomerID    int       `json:"customer_id"`
	ProductID     int       `json:"product_id"`
	TransactionID int       `json:"transaction_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// InsertOrder inserts a new order, and returns its id
func (dbm *DBModels) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into orders 
			(product_id, transaction_id, status_id, quantity, customer_id, amount, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := dbm.DB.ExecContext(ctx, stmt,
		order.ProductID,
		order.TransactionID,
		order.StatusID,
		order.Quantity,
		order.CustomerID,
		order.Amount,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
