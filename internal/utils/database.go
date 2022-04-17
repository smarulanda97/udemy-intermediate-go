package utils

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConfig struct {
	Dsn string
}

func OpenDBConnection(dsn string) (*sql.DB, error) {
	connection, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := connection.Ping(); err != nil {
		return nil, err
	}

	return connection, nil
}
