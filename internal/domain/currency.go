package domain

import (
    "golang.org/x/text/language"
    "golang.org/x/text/message"
)

// Format formats a float64 as a currency string with two decimals.
func Format(v float64) string {
    // Use a printer for English locale to format with commas and a leading $ sign.
    p := message.NewPrinter(language.English)
    return p.Sprintf("$%0.2f", v)
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
