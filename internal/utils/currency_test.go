package utils

import (
	"encoding/json"
	"testing"
)

func TestCurrencyValidate(t *testing.T) {
	cases := []struct {
		cur     Currency
		wantErr bool
	}{
		{"USD", false},
		{"usd", true},
		{"EU", true},
		{"EURO", true},
		{"", true},
	}
	for _, c := range cases {
		err := c.cur.Validate()
		if (err != nil) != c.wantErr {
			t.Fatalf("Validate(%s) error = %v, wantErr %v", c.cur, err, c.wantErr)
		}
	}
}

func TestCurrencyJSON(t *testing.T) {
	var cur Currency = "EUR"
	data, err := json.Marshal(cur)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	var cur2 Currency
	if err := json.Unmarshal(data, &cur2); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if cur2 != cur {
		t.Fatalf("expected %s got %s", cur, cur2)
	}
	// invalid JSON should fail
	var bad Currency
	if err := json.Unmarshal([]byte(`"usd"`), &bad); err == nil {
		t.Fatalf("expected error for lowercase currency")
	}
}
