package db

import "database/sql"

// DB is a placeholder interface for database operations.
// Future methods will be added as the application grows.
type DB interface {}

type db struct {
    conn *sql.DB
}

// NewDB creates a new DB implementation wrapping the given sql.DB connection.
func NewDB(conn *sql.DB) DB {
    return &db{conn: conn}
}
