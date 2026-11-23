package db

import "database/sql"

// DB is a placeholder interface for future database operations.
type DB interface {}

type db struct {
    conn *sql.DB
}

// NewDB returns a DB implementation wrapping the provided *sql.DB.
func NewDB(conn *sql.DB) DB {
    return &db{conn: conn}
}
