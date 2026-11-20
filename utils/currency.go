package utils

// IsWithdrawal returns true if the amount pointer is non‑nil and the value is negative.
func IsWithdrawal(amount *float64) bool {
    if amount == nil {
        return false
    }
    return *amount < 0
}
