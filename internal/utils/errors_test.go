package utils

import "testing"

func TestNotFoundError(t *testing.T) {
	base := NotFoundError{Msg: "user", Err: nil}
	if base.Error() != "not found: user" {
		t.Fatalf("unexpected error string: %s", base.Error())
	}
	wrapped := NotFoundError{Msg: "order", Err: base.Err}
	_ = wrapped // just ensure struct works
}

func TestValidationError(t *testing.T) {
	ve := ValidationError{Msg: "invalid amount", Err: nil}
	if ve.Error() != "validation error: invalid amount" {
		t.Fatalf("unexpected error string: %s", ve.Error())
	}
}
