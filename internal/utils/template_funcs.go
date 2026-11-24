package utils

import "fmt"

// FuncMap returns a map of helper functions for html/template.
func FuncMap() map[string]interface{} {
    return map[string]interface{}{
        "IsWithdrawal": IsWithdrawal,
        "CurrencyFromDouble": CurrencyFromDouble,
    }
}

// IsWithdrawal returns true if the amount is negative (a withdrawal).
func IsWithdrawal(amount float64) bool { return amount < 0 }

// CurrencyFromDouble formats a float as a currency string with two decimals.
func CurrencyFromDouble(v float64) string { return fmt.Sprintf("%.2f", v) }
