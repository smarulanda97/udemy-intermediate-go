package models

import (
	"context"
	"time"
)

// Product is the type for all products
type Product struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// GetProduct returns a product by its id
func (m *DBModels) GetProduct(id int) (Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var p Product

	stmt := `
		select 
			id, name, description, price, inventory_level, coalesce(image, ''), is_recurring, plan_id, created_at, updated_at
		from 
			products 
		where id = ?
	`

	row := m.DB.QueryRowContext(ctx, stmt, id)

	err := row.Scan(&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.InventoryLevel,
		&p.Image,
		&p.IsRecurring,
		&p.PlanID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		return p, nil
	}

	return p, nil
}
