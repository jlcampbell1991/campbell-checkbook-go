package db

import "database/sql"

type DB interface{}

type db struct {
	conn *sql.DB
}

// NewDB creates a new DB implementation wrapping the provided *sql.DB.
func NewDB(conn *sql.DB) DB {
	return &db{conn: conn}
}
