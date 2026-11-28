package domain

import "fmt"

// Format formats a float64 as a currency string with two decimals.
func Format(v float64) string {
	return fmt.Sprintf("%.2f", v)
}

// Value returns the dereferenced float64 or zero if nil.
func Value(v *float64) float64 {
    if v == nil {
        return 0
    }
    return *v
}

// IsWithdrawal returns the CSS class name ("withdrawal" or "deposit") based on sign.
func IsWithdrawal(v float64) string {

	if v < 0 {
		return "withdrawal"
	}
	return "deposit"
}
