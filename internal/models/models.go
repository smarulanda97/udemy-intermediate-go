package models

import (
	"database/sql"
)

type DBModels struct {
	DB *sql.DB
}

type Models struct {
	DB DBModels
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModels{
			DB: db,
		},
	}
}

// SaveCustomer saves customer and returns id
func (dbm *DBModels) SaveCustomer(customer Customer) (int, error) {
	id, err := dbm.InsertCustomer(customer)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SaveTransaction saves customer and returns id
func (dbm *DBModels) SaveTransaction(txn Transaction) (int, error) {
	id, err := dbm.InsertTransaction(txn)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// SaveOrder saves customer and returns id
func (dbm *DBModels) SaveOrder(order Order) (int, error) {
	id, err := dbm.InsertOrder(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}
