package utils

import (
	"database/sql"
	"errors"
)

// ScanNullableString scans a nullable SQL string into a *string. If the column is NULL the destination is set to nil.
func ScanNullableString(src interface{}, dest **string) error {
	if src == nil {
		*dest = nil
		return nil
	}
	// Expect src to be []byte or string.
	switch v := src.(type) {
	case []byte:
		s := string(v)
		*dest = &s
		return nil
	case string:
		s := v
		*dest = &s
		return nil
	default:
		return errors.New("unsupported type for ScanNullableString")
	}
}

// ExecWithTx runs the supplied function within a transaction. It commits on success and rolls back on error.
func ExecWithTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		// Ensure rollback if panic or error not handled.
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()
	if err = fn(tx); err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return errors.Join(err, rbErr)
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
