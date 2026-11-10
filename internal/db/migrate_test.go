package db

import "testing"

func TestRunMigrations_NoDSN(t *testing.T) {
    // Should not panic when DSN is empty.
    RunMigrations("")
    // If we reach here, the function handled empty DSN gracefully.
}
