package utils

import "testing"

func TestIsWithdrawal(t *testing.T) {
	if !IsWithdrawal(-5) {
		t.Errorf("expected true for negative amount")
	}
	if IsWithdrawal(0) {
		t.Errorf("expected false for zero amount")
	}
	if IsWithdrawal(10) {
		t.Errorf("expected false for positive amount")
	}
}

func TestCurrencyFromDouble(t *testing.T) {
	got := CurrencyFromDouble(12.3)
	if got != "12.30" {
		t.Errorf("expected 12.30, got %s", got)
	}
}
