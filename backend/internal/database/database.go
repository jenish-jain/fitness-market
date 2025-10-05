package database

import (
	"database/sql"
	"fmt"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}