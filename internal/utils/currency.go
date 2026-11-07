package utils

import (
	"encoding/json"
	"errors"
	"strings"
)

// Currency represents an ISO 4217 currency code (e.g., "USD", "EUR").
type Currency string

// Validate checks that the currency code is non‑empty, three uppercase letters.
func (c Currency) Validate() error {
	s := string(c)
	if len(s) != 3 {
		return errors.New("currency must be three characters")
	}
	if s != strings.ToUpper(s) {
		return errors.New("currency must be uppercase")
	}
	return nil
}

// MarshalJSON implements json.Marshaler. It validates before marshaling.
func (c Currency) MarshalJSON() ([]byte, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(string(c))
}

// UnmarshalJSON implements json.Unmarshaler. It validates after unmarshaling.
func (c *Currency) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	cur := Currency(s)
	if err := cur.Validate(); err != nil {
		return err
	}
	*c = cur
	return nil
}
