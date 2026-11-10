package db

import (
    "database/sql"
    "fmt"
    "log"
)

// InsertAccount inserts a new account with the given name and returns its ID.
func InsertAccount(db *sql.DB, name string) (int64, error) {
    if name == "" {
        return 0, fmt.Errorf("account name cannot be empty")
    }
    var id int64
    err := db.QueryRow("INSERT INTO accounts (name) VALUES ($1) RETURNING id", name).Scan(&id)
    if err != nil {
        return 0, err
    }
    log.Printf("created account %d with name %s", id, name)
    return id, nil
}
